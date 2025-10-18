package payment

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/darkphotonKN/stripe-advanced-approach/internal/interfaces"
	"github.com/darkphotonKN/stripe-advanced-approach/internal/model"
	"github.com/darkphotonKN/stripe-advanced-approach/internal/user"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	redislib "github.com/redis/go-redis/v9"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82/subscription"
)

type service struct {
	userService      PaymentUserService
	paymentProcessor PaymentProcessor
	cacheClient      interfaces.Cache
	repo             Repository
}

type Repository interface {
	Create(ctx context.Context, userId uuid.UUID, paymentIntent *PaymentIntentRequest) error
	GetPaymentByIntentID(ctx context.Context, intentID string) (*Payment, error)
	UpdateStatus(ctx context.Context, intentID string, status string) error
	UpsertPayment(ctx context.Context, paymentIntentID string, payment *Payment) error
	UpsertSubscriptionRecord(ctx context.Context, sub *Subscription) error
	GetActiveSubscription(ctx context.Context, userID uuid.UUID) (*Subscription, error)
	UpdateSubscriptionStatus(ctx context.Context, subID string, status string) error
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
}

type PaymentUserService interface {
	UpdateStripeCustomer(ctx context.Context, userID uuid.UUID, stripeCustomerID string) error
	GetByStripeCustomerID(ctx context.Context, stripeCustomerID string) (*user.User, error)
}

func NewService(repo Repository, userService PaymentUserService, paymentProcessor PaymentProcessor, cacheClient interfaces.Cache) *service {
	return &service{
		repo:             repo,
		userService:      userService,
		paymentProcessor: paymentProcessor,
		cacheClient:      cacheClient,
	}
}

/*
*
*
* Primary method for syncing up stripe-related states and avoiding a split-brain problem.
*
* Since payment table's status comes from stripe, we need to at least get that validated from stripe's more
* consistent apis, as opposed to webhooks, to store in the KV for easy access but at the same time update the database
* once we have it validated.
*
* The key-value cache structure will be as follows:

Key: stripe:customer:{customerId}

	Value: {
	  subscriptionId: "sub_xyz",
	  status: "active",           // Core validation field
	  priceId: "price_abc",       // Feature access control
	  currentPeriodEnd: 1234567,  // Billing cycle info
	  cancelAtPeriodEnd: false    // Immediate cancellation status
	}

Key: stripe:user:{userId}
Value: "cus_stripe123"        // Customer ID mapping

* Database storage will store the same things but into the respective areas
*
*/
func (s *service) SyncStripeDataToStorage(ctx context.Context, customerId string) error {
	// --- Data Organization ---

	// get latest up-to-date data from stripe
	customer, err := customer.Get(customerId, nil)

	if err != nil {
		return fmt.Errorf("failed to get customer from stripe: %w", err)
	}

	customerJSON, err := json.MarshalIndent(customer, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal customer data: %w", err)
	}

	fmt.Printf("\n=== Stripe Customer Data ===\n%s\n============================\n\n", string(customerJSON))

	stripeCusKey := fmt.Sprintf("stripe:customer:%s", customerId)

	// -- subscriptions --

	params := &stripe.SubscriptionListParams{
		Customer: stripe.String(customerId),
		Status:   stripe.String("all"),
	}

	// expand the payment method to get card details
	params.AddExpand("data.default_payment_method")

	// get subscription data
	subIter := subscription.List(params)

	// subscriptions slice
	subscriptions := []*stripe.Subscription{}

	// validates that customer has subscriptions
	for subIter.Next() {
		// get the subscription
		sub := subIter.Subscription()

		// Handle iteration error
		if err := subIter.Err(); err != nil {
			fmt.Printf("failed to fetch subscriptions from Stripe: %+v", err)
			return fmt.Errorf("failed to fetch subscriptions from Stripe: %w", err)
		}

		subscriptions = append(subscriptions, sub)
	}

	// -- payments --
	payments := []*stripe.PaymentIntent{}

	paymentParams := &stripe.PaymentIntentListParams{
		Customer: stripe.String(customerId),
	}

	// include payment method details
	paymentParams.AddExpand("data.payment_method")

	paymentIter := paymentintent.List(paymentParams)

	for paymentIter.Next() {
		pi := paymentIter.PaymentIntent()
		fmt.Printf("\npayment intent: %+v\n\n", pi)

		payments = append(payments, pi)
	}

	// --- DB Storage ---

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("error when attempting to start transaction: %+v", err)
	}

	// NOTE: safe to run even if commit was successful - in that case it will be a no-op
	defer tx.Rollback()

	// update application database for the respective tables

	// -- user
	user, err := s.userService.GetByStripeCustomerID(ctx, customerId)

	// -- subscription --

	for _, sub := range subscriptions {
		err := s.repo.UpsertSubscriptionRecord(ctx, &Subscription{
			UserID:               user.ID,
			Status:               string(sub.Status),
			StripeSubscriptionID: sub.ID,
			StripeCustomerID:     customerId,
		})

		if err != nil {
			return fmt.Errorf("\nError when attempting to batch upsert subscriptions during sync: %+v\n\n", err)
		}
	}

	// -- payments --

	// get userId from cache / db depending on availability
	userId, err := s.GetCachedUserIdByCustomerId(ctx, customerId)

	for _, payment := range payments {
		err := s.repo.UpsertPayment(ctx, payment.ID, &Payment{
			UserID:           userId,
			StripeCustomerID: customerId,
			StripeIntentID:   payment.ID,
			Amount:           payment.Amount,
			Status:           string(payment.Status),
			Currency:         string(payment.Currency),
		})

		if err != nil {
			return err
		}
	}

	// --- Caching ---

	subCache := make([]*StripeSubscriptionCache, len(subscriptions))

	// -- subscriptions --

	for index, sub := range subscriptions {
		// extract payment method info safely
		var pmInfo *PaymentMethodInfo

		if len(subscriptions) > 0 {
			sub := subscriptions[0]

			if sub.DefaultPaymentMethod != nil && sub.DefaultPaymentMethod.Card != nil {
				pmInfo = &PaymentMethodInfo{
					Brand: string(sub.DefaultPaymentMethod.Card.Brand),
					Last4: sub.DefaultPaymentMethod.Card.Last4,
				}
			}

		}

		// add to cache slice
		subCache[index] = &StripeSubscriptionCache{
			SubscriptionID:    sub.ID,
			Status:            string(sub.Status),
			PriceID:           sub.Items.Data[0].Price.ID,
			CancelAtPeriodEnd: sub.CancelAtPeriodEnd,
			PaymentMethod:     pmInfo,
		}
	}

	// -- payments --

	paymentCache := make([]*StripePaymentsCache, len(payments))

	for index, payment := range payments {
		paymentCache[index] = &StripePaymentsCache{
			ID:     payment.ID,
			Status: string(payment.Status),
		}
	}

	// -- user --

	stripeCusData := StripeCustomerDataRes{
		ID:                   customer.ID,
		Address:              convertAddress(customer.Address),
		Balance:              customer.Balance,
		CashBalance:          convertCashBalance(customer.CashBalance),
		Created:              customer.Created,
		Currency:             string(customer.Currency),
		DefaultSource:        convertDefaultSource(customer.DefaultSource),
		Deleted:              customer.Deleted,
		Delinquent:           customer.Delinquent,
		Description:          customer.Description,
		Discount:             convertDiscount(customer.Discount),
		Email:                customer.Email,
		InvoiceCreditBalance: customer.InvoiceCreditBalance,
		InvoicePrefix:        customer.InvoicePrefix,
		InvoiceSettings:      convertInvoiceSettings(customer.InvoiceSettings),
		Livemode:             customer.Livemode,
		Metadata:             customer.Metadata,
		Name:                 customer.Name,
		NextInvoiceSequence:  customer.NextInvoiceSequence,
		Object:               customer.Object,
		Phone:                customer.Phone,
		PreferredLocales:     customer.PreferredLocales,
		Subscriptions:        convertSubscriptions(customer.Subscriptions),
		Tax:                  convertTax(customer.Tax),
		TaxExempt:            string(customer.TaxExempt),
	}

	// combine the two pieces of information into one cache state
	cacheState := StripeCacheData{
		CustomerData:  stripeCusData,
		Subscriptions: subCache,
		Payments:      paymentCache,
	}

	cacheStateJSON, err := json.Marshal(cacheState)

	if err != nil {
		return fmt.Errorf("failed to marshal cacheState: %w", err)
	}

	// update redis
	err = s.cacheClient.Set(ctx, stripeCusKey, cacheStateJSON, 0)

	if err != nil {
		return fmt.Errorf("failed to sync and store stripe data into cache: %w", err)
	}

	return nil
}

/**
* adds/sets the mapping between userId and payment processor customerId in cache
**/
func (s *service) AddCacheUserIdToCusId(ctx context.Context, userId uuid.UUID, customerId string) error {
	key := fmt.Sprintf("stripe:userId:%s:customerId", userId.String())
	err := s.cacheClient.Set(ctx, key, customerId, 0)

	if err != nil {
		return fmt.Errorf("failed to cache userId to customerId mapping: %w", err)
	}

	return nil
}

/**
* gets cached customerId with the userId
**/
func (s *service) GetCachedCusIdFromUserId(ctx context.Context, userId uuid.UUID) (string, error) {
	key := fmt.Sprintf("stripe:userId:%s:customerId", userId.String())
	customerId, err := s.cacheClient.Get(ctx, key)

	// doesn't exist in cache, error, cache is supposed to have a mapping from this point
	if err == redislib.Nil {
		return "", fmt.Errorf("No customerId exists for this userId %s", userId)
	}

	if err != nil {
		return "", fmt.Errorf("unexpected error occured when attempting to find map of customer Id from userId: %s", userId)
	}

	return customerId, nil
}

/**
* The get version of the payment cache sync method. Gets the latest up-to-date data from the cache if it exists,
* otherwise calls the sync method to update the cache.
**/
func (s *service) GetStripeData(ctx context.Context, customerId string) (*StripeCacheData, error) {
	// check if customer data already exists in the cache
	dataJSON, err := s.cacheClient.Get(ctx, customerId)

	// if it doesn't we sync the data right there
	if err == redislib.Nil {
		fmt.Printf("Customer data doesn't exist in cache.")

		// sync data to cache
		err := s.SyncStripeDataToStorage(ctx, customerId)
		if err != nil {
			fmt.Printf("failed to resync method during attempt to get cached data.")
			return nil, err
		}

		// attempt to get the data again after syncing
		return s.GetStripeData(ctx, customerId)
	}

	// handle other exceptions
	if err != nil {
		fmt.Printf("error when attempting to get cache data for customerID %s\nerr was:\n%+v\n", customerId, err)
		return nil, err
	}

	// data already exists, just unmarshal and return it
	fmt.Printf("\ncache data before unmarshal: %+v\n\n", dataJSON)

	var cacheData StripeCacheData

	err = json.Unmarshal([]byte(dataJSON), &cacheData)
	if err != nil {
		fmt.Printf("error when attempting to get unmarshal data for customerID %s\nerr was:\n%+v\n", customerId, err)
		return nil, err
	}

	fmt.Printf("\ncache data after unmarshal: %+v\n\n", cacheData)
	return &cacheData, nil
}

func (s *service) SetupProducts(ctx context.Context, request *SetupProductsReq) (*SetupProductsResp, error) {
	return s.paymentProcessor.SetupProducts(ctx, request)
}

func (s *service) CreateCustomer(ctx context.Context, userId uuid.UUID, email string) (string, error) {
	// create customer on stripe and get customer id
	customerId, err := s.paymentProcessor.CreateCustomer(ctx, userId, email)

	if err != nil {
		fmt.Printf("Error occured when attemtping to create customer on stripe, %s\n", err.Error())
		return "", err
	}

	// update local user repo for mapping
	err = s.userService.UpdateStripeCustomer(ctx, userId, customerId)

	if err != nil {
		fmt.Printf("Error occured when attempting to update stripe customerId to user repo in CreateCustomer method: %s\n", err.Error())
		return "", err
	}

	return customerId, nil
}

func (s *service) SaveCard(ctx context.Context, customerId string) (string, error) {
	return s.paymentProcessor.SaveCard(ctx, customerId)
}

func (s *service) CreatePaymentIntent(ctx context.Context, amount int64, customerId string) (*CreatePaymentIntentResponse, error) {
	return s.paymentProcessor.CreatePaymentIntent(ctx, amount, customerId)
}

func (s *service) GetProducts(ctx context.Context) (*ProductListResponse, error) {
	return s.paymentProcessor.GetProducts(ctx)
}

func (s *service) PurchaseProduct(ctx context.Context, userId uuid.UUID, req *PurchaseProductRequest) (*PurchaseProductResponse, error) {
	res, err := s.paymentProcessor.PurchaseProduct(ctx, req)

	if err != nil {
		return nil, err
	}

	// create payments record in database to map payment status to that on the payment service

	fmt.Printf("\npurchase product request: %+v\n\n", req)
	fmt.Printf("\npurchase product payment processor response: %+v\n\n", res)

	err = s.repo.Create(ctx, userId, &PaymentIntentRequest{
		CustomerID: req.CustomerID,
		Amount:     res.Amount,
		IntentID:   res.PaymentIntentID,
	})

	if err != nil {
		return nil, err
	}

	return &PurchaseProductResponse{
		ClientSecret:    res.ClientSecret,
		PaymentIntentID: res.PaymentIntentID,
	}, nil
}

func (s *service) SetupSubscription(ctx context.Context, request *SetupProductsReq) (*SetupProductsResp, error) {
	return s.paymentProcessor.SetupSubscription(ctx, request)
}

/**
* Default Stripe-Recommended Flow:
*
* 1. Subscribe action
* User clicks "Subscribe" →  Backend creates Stripe subscription with
* `payment_behavior: "default_incomplete"` →  Returns client_secret
* → Frontend confirms payment with Stripe Elements →  User completes payment →
* Stripe redirects to success page
*
* 2. Database Storage (naive implementation, but recommended by Stripe)
* When subscription created → Store in DB as status: "incomplete" →  Wait for webhooks to update status to "active"
**/
func (s *service) SubscribeToProduct(ctx context.Context, userId uuid.UUID, req *SubscribeRequest) (*SubscribeResponse, error) {
	res, err := s.paymentProcessor.SubscribeToProduct(ctx, req)

	if err != nil {
		return nil, err
	}

	// store in database a new subscription with data from successful Stripe response
	err = s.repo.UpsertSubscriptionRecord(ctx, &Subscription{
		UserID:               userId,
		StripeCustomerID:     req.CustomerID,
		StripeSubscriptionID: res.SubscriptionID,
		Status:               res.Status,
	})

	if err != nil {
		fmt.Printf("\nError when creating a subscription record in DB: %+v\n\n", err)
		return nil, err
	}

	return res, nil
}

func (s *service) GetCachedUserIdByCustomerId(ctx context.Context, customerID string) (uuid.UUID, error) {
	key := fmt.Sprintf("stripe:customer:%s:userid", customerID)
	userIdStr, err := s.cacheClient.Get(ctx, key)

	// key doesn't exist, acquire userId to fill in cache
	if err == redislib.Nil {
		user, err := s.userService.GetByStripeCustomerID(ctx, customerID)

		if err != nil {
			fmt.Printf("err when attempting to get user with customerId %s: %v\n", customerID, err)
			return uuid.Nil, err
		}

		// store in cache
		s.cacheClient.Set(ctx, key, user.ID.String(), 0)

		// return the id
		return user.ID, nil
	}

	// unexpected errors
	if err != nil {
		fmt.Printf("err when attempting to get userId from cache with customerId %s: %v\n", customerID, err)
		return uuid.Nil, err
	}

	// convert back to uuid
	userId, err := uuid.Parse(userIdStr)

	if err != nil {
		return uuid.Nil, err
	}

	return userId, err
}

// --- Full Flow Methods ---

/**
* Recieves a payment processor event and parses the event
**/
func (s *service) ProcessWebhookEvent(ctx context.Context, event *stripe.Event) error {
	customerId, err := s.paymentProcessor.ProcessWebhookEvent(ctx, event)

	if err != nil {
		fmt.Printf("\npaymentProcessor method ProcessWebhookEvent could not process incoming event of %s, err :+v\n\n", event, err)
		return err
	}

	fmt.Printf("Service layer - customerId: %s\n", customerId)

	s.SyncStripeDataToStorage(ctx, customerId)

	return nil
}

// Helper functions to convert Stripe types to our cache types
func convertAddress(addr *stripe.Address) *CustomerAddress {
	if addr == nil {
		return nil
	}
	return &CustomerAddress{
		City:       addr.City,
		Country:    addr.Country,
		Line1:      addr.Line1,
		Line2:      addr.Line2,
		PostalCode: addr.PostalCode,
		State:      addr.State,
	}
}

func convertCashBalance(cb *stripe.CashBalance) *CustomerCashBalance {
	if cb == nil {
		return nil
	}
	return &CustomerCashBalance{
		Object:    cb.Object,
		Available: cb.Available,
		Customer:  cb.Customer,
		Livemode:  cb.Livemode,
	}
}

func convertDefaultSource(ds *stripe.PaymentSource) *string {
	if ds == nil {
		return nil
	}
	id := ds.ID
	return &id
}

func convertDiscount(d *stripe.Discount) *CustomerDiscount {
	if d == nil {
		return nil
	}
	return &CustomerDiscount{
		Coupon:   convertCoupon(d.Coupon),
		Customer: d.ID,
		End:      d.End,
		Id:       d.ID,
		Object:   d.Object,
		Start:    d.Start,
	}
}

func convertCoupon(c *stripe.Coupon) *Coupon {
	if c == nil {
		return nil
	}
	return &Coupon{
		Id:               c.ID,
		Object:           c.Object,
		AmountOff:        c.AmountOff,
		Created:          c.Created,
		Currency:         string(c.Currency),
		Duration:         string(c.Duration),
		DurationInMonths: c.DurationInMonths,
		Livemode:         c.Livemode,
		MaxRedemptions:   c.MaxRedemptions,
		Name:             c.Name,
		PercentOff:       c.PercentOff,
		RedeemBy:         c.RedeemBy,
		TimesRedeemed:    c.TimesRedeemed,
		Valid:            c.Valid,
	}
}

func convertInvoiceSettings(is *stripe.CustomerInvoiceSettings) *CustomerInvoiceSettings {
	if is == nil {
		return nil
	}

	var customFields []*InvoiceCustomField
	for _, cf := range is.CustomFields {
		if cf != nil {
			customFields = append(customFields, &InvoiceCustomField{
				Name:  cf.Name,
				Value: cf.Value,
			})
		}
	}

	var defaultPM *string
	if is.DefaultPaymentMethod != nil {
		id := is.DefaultPaymentMethod.ID
		defaultPM = &id
	}

	var renderingOpts *InvoiceRenderingOptions
	if is.RenderingOptions != nil {
		renderingOpts = &InvoiceRenderingOptions{
			AmountTaxDisplay: string(is.RenderingOptions.AmountTaxDisplay),
		}
	}

	return &CustomerInvoiceSettings{
		CustomFields:         customFields,
		DefaultPaymentMethod: defaultPM,
		Footer:               is.Footer,
		RenderingOptions:     renderingOpts,
	}
}

func convertSubscriptions(s *stripe.SubscriptionList) *SubscriptionList {
	if s == nil {
		return nil
	}

	var data []interface{}
	for _, sub := range s.Data {
		data = append(data, sub)
	}

	return &SubscriptionList{
		Data:    data,
		HasMore: s.HasMore,
		Url:     s.URL,
	}
}

func convertTax(t *stripe.CustomerTax) *CustomerTax {
	if t == nil {
		return nil
	}

	var location *CustomerTaxLocation
	if t.Location != nil {
		location = &CustomerTaxLocation{
			Country: t.Location.Country,
			Source:  string(t.Location.Source),
			State:   t.Location.State,
		}
	}

	return &CustomerTax{
		AutomaticTax: string(t.AutomaticTax),
		IpAddress:    t.IPAddress,
		Location:     location,
	}
}

func convertTaxIds(t *stripe.TaxIDList) *CustomerTaxIdList {
	if t == nil {
		return nil
	}

	var data []CustomerTaxId
	for _, taxId := range t.Data {
		if taxId != nil {
			data = append(data, CustomerTaxId{
				Id:      taxId.ID,
				Object:  taxId.Object,
				Country: taxId.Country,
				Type:    string(taxId.Type),
				Value:   taxId.Value,
			})
		}
	}

	return &CustomerTaxIdList{
		Data:    data,
		HasMore: t.HasMore,
		Url:     t.URL,
	}
}

/**
* for utilizing cache for checking the user's subscription status to the pro
* plan of this site
**/
func (s *service) GetSubscriptionStatusCache(ctx context.Context, userId uuid.UUID) (*model.SubscriptionStatus, error) {
	// get customerId from cache
	cusId, err := s.GetCachedCusIdFromUserId(ctx, userId)

	if err != nil {
		return nil, err
	}

	// TODO:
	// get subscription status
	stripeCacheData, err := s.GetStripeData(ctx, cusId)

	fmt.Printf("stripeCachData when getting subscription status:", stripeCacheData)
	return nil, nil
}

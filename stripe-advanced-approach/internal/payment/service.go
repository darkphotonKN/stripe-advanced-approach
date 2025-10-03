package payment

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/darkphotonKN/stripe-advanced-approach/internal/interfaces"
	"github.com/google/uuid"
	redislib "github.com/redis/go-redis/v9"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
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
	CreateSubscriptionRecord(ctx context.Context, sub *Subscription) error
	GetActiveSubscription(ctx context.Context, userID uuid.UUID) (*Subscription, error)
	UpdateSubscriptionStatus(ctx context.Context, subID string, status string) error
}

type PaymentUserService interface {
	UpdateStripeCustomer(ctx context.Context, userID uuid.UUID, stripeCustomerID string) error
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
	stripeCusKey := fmt.Sprintf("stripe:customer:%s", customerId)

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

	// get subscription data
	params := &stripe.SubscriptionListParams{
		Customer: stripe.String(customerId),
		Status:   stripe.String("all"),
	}

	// expand the payment method to get card details
	params.AddExpand("data.default_payment_method")

	iter := subscription.List(params)

	// validates that customer has subscriptions
	if !iter.Next() {
		// no subscription
		noSubData := map[string]interface{}{
			"status": "none",
		}

		noSubJSON, _ := json.Marshal(noSubData)

		err := s.cacheClient.Set(ctx, stripeCusKey, noSubJSON, 0)
		if err != nil {
			return fmt.Errorf("failed to cache no subscription data: %w", err)
		}

		return nil
	}

	// Get the subscription
	sub := iter.Subscription()

	// Handle iteration error
	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to fetch subscriptions from Stripe: %w", err)
	}

	// extract payment method info safely
	var pmInfo *PaymentMethodInfo
	if sub.DefaultPaymentMethod != nil && sub.DefaultPaymentMethod.Card != nil {
		pmInfo = &PaymentMethodInfo{
			Brand: string(sub.DefaultPaymentMethod.Card.Brand),
			Last4: sub.DefaultPaymentMethod.Card.Last4,
		}
	}

	// build the subscriptions / payment cache data
	subCache := StripeSubscriptionCache{
		SubscriptionID:    sub.ID,
		Status:            string(sub.Status),
		PriceID:           sub.Items.Data[0].Price.ID,
		CancelAtPeriodEnd: sub.CancelAtPeriodEnd,
		PaymentMethod:     pmInfo,
	}

	// customer data
	stripeCusData := json.Unmarshal()

	// combine the two pieces of information into one cache state
	cacheState := StripeCacheData{
		customerData: StripeCustomerDataRes{
			ID: customer.ID,
		},
		subscriptions: subCache,
	}

	// update redis
	err = s.cacheClient.Set(ctx, stripeCusKey, customerJSON, 0)

	if err != nil {
		return fmt.Errorf("failed to sync and store stripe data into cache: %w", err)
	}

	// update application database for the respective tables

	// user

	// payments

	return nil
}

/**
* The get version of the stripe sync method. Gets the latest up-to-date data from the cache if it exists,
* otherwise calls the sync method to update the cache.
**/
func (s *service) GetStripeData(ctx context.Context, customerId string) (*StripeCacheData, error) {
	// check if customer data already exists in the cache
	data, err := s.cacheClient.Get(ctx, customerId)

	// if it doesn't we sync the data right there
	if err != redislib.Nil {
		fmt.Printf("Customer data doesn't exist in cache.")

		// sync data to cache
		s.SyncStripeDataToStorage(ctx, customerId)
	}

	// handle other exceptions
	if err != nil {
		fmt.Printf("error when attempting to get cache data for customerID %s\nerr was:\n%+v\n", customerId, err)
		return nil, err
	}

	// data already exists, just unmarshal and return it
	fmt.Printf("\ncache data before unmarshal: %+v\n\n", data)

	var cacheData StripeCacheData

	err = json.Unmarshal([]byte(data), &cacheData)

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
	err = s.repo.CreateSubscriptionRecord(ctx, &Subscription{
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

// --- Full Flow Methods ---

/**
* Recieves a stripe event and parses the event
**/
func (s *service) ProcessWebhookEvent(ctx context.Context, event *stripe.Event) error {
	fmt.Printf("Processing webhook event type: %s\n", event.Type)

	parsedPaymentIntent, err := s.parsePaymentProcessorEvent(event)

	if err != nil {
		return err
	}

	fmt.Printf("Processing webhook event type: %s\n", event.Type)

	switch event.Type {
	case "payment_intent.succeeded":
		return s.repo.UpdateStatus(ctx, parsedPaymentIntent.ID, "success")
	case "payment_intent.payment_failed":
		return s.repo.UpdateStatus(ctx, parsedPaymentIntent.ID, "failed")
	case "payment_intent.canceled":
		return s.repo.UpdateStatus(ctx, parsedPaymentIntent.ID, "canceled")

	// TODO: add other subscription statuses
	case "customer.subscription.created":
		return s.repo.UpdateSubscriptionStatus(ctx, parsedPaymentIntent.ID, "active")
	default:
		fmt.Printf("Unhandled event type: %s\n", event.Type)
		return nil
	}
}

func (s *service) parsePaymentProcessorEvent(event *stripe.Event) (*stripe.PaymentIntent, error) {
	var paymentIntent stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &paymentIntent); err != nil {
		return nil, fmt.Errorf("error parsing payment intent: %w", err)
	}

	return &paymentIntent, nil
}

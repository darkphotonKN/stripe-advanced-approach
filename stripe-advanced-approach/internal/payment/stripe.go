package payment

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82/price"
	"github.com/stripe/stripe-go/v82/product"
	"github.com/stripe/stripe-go/v82/setupintent"
	"github.com/stripe/stripe-go/v82/subscription"
)

type StripeProcessor struct{}

func NewStripeProcessor() PaymentProcessor {
	return &StripeProcessor{}
}

/**
* Sets up custom products for purchase. This can be anything you wish to sell, and is a digital representation of the
* item - which can be physical, digital, or even just a concept (donation, etc).
**/
func (s *StripeProcessor) SetupProducts(ctx context.Context, request *SetupProductsReq) (*SetupProductsResp, error) {
	// create product
	prod, err := product.New(&stripe.ProductParams{
		Name:        stripe.String(request.Name),
		Description: stripe.String(request.Description),
	})

	if err != nil {
		fmt.Printf("\nError when creating product on stripe: %+v\n\n", err)
		return nil, err
	}

	// create STANDARD product / service
	oneTimePrice, err := price.New(&stripe.PriceParams{
		Currency: stripe.String("usd"),
		Product:  stripe.String(prod.ID),
		// NO Recurring parameter = one-time price!
		UnitAmount: stripe.Int64(request.Price),
	})

	if err != nil {
		fmt.Printf("\nError when creating subscription price for product on stripe: %+v\n\n", err)
		return nil, err
	}

	fmt.Printf("Created new product's price successfully. Response:%+v\n", oneTimePrice)

	// set default price. NOT set by default.
	product.Update(prod.ID, &stripe.ProductParams{
		DefaultPrice: stripe.String(oneTimePrice.ID),
	})

	return &SetupProductsResp{
		PriceID: oneTimePrice.ID,
	}, nil
}

func (s *StripeProcessor) CreateCustomer(ctx context.Context, userId uuid.UUID, email string) (string, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
		Metadata: map[string]string{
			"user_id": userId.String(),
		},
	}

	// create customer on Stripe
	cust, err := customer.New(params)
	if err != nil {
		return "", err
	}

	fmt.Printf("Created customer: %s", cust.ID)

	return cust.ID, nil
}

/**
* This method AUTHORIZES a card save for a customer by creating a permission token for the client
* to then use the stripe sdk via elements to save the card.
**/
func (s *StripeProcessor) SaveCard(ctx context.Context, customerId string) (string, error) {
	params := &stripe.SetupIntentParams{
		Customer: stripe.String(customerId),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
	}

	// - creates and sets up authorization for FUTURE card purchases.
	// - generates client_secret, a permission token. NO card data is saved.
	// - links to customer in stripe's system via customerId
	si, err := setupintent.New(params)
	if err != nil {
		fmt.Printf("Error when attempting to generate setup intent: %s\n", err.Error())
		return "", err
	}

	return si.ClientSecret, nil
}

/**
* Creates a payment authorization token that allows the frontend to charge a specific
* amount for a specific customer. The backend validates the request and gets Stripe's
* permission to charge, but the actual payment happens when the frontend confirms
* with card data. This prevents unauthorized charges while keeping card data secure.
**/
func (s *StripeProcessor) CreatePaymentIntent(ctx context.Context, amount int64, customerId string) (*CreatePaymentIntentResponse, error) {

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String("usd"),
		Customer: stripe.String(customerId),

		// manual confirmation means frontend confirms - this is default confirmation
		ConfirmationMethod: stripe.String("automatic"),
		Confirm:            stripe.Bool(false),

		// only allow CARD
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
	}

	intent, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	fmt.Printf("stripe created intent: %+v\n", intent)

	return &CreatePaymentIntentResponse{
		PaymentIntentID: intent.ID,
		ClientSecret:    intent.ClientSecret,
	}, nil
}

/**
* Creates a subscription item or service for recurring type payments.
**/
func (s *StripeProcessor) SetupSubscription(ctx context.Context, request *SetupProductsReq) (*SetupProductsResp, error) {
	// create subscription
	subscriptionProd, err := product.New(&stripe.ProductParams{
		Name:        stripe.String(request.Name),
		Description: stripe.String(request.Description),
	})

	if err != nil {
		fmt.Printf("\nError when creating product on stripe: %+v\n\n", err)
		return nil, err
	}

	// create SUBSCRIPTION product / service
	subscriptionPrice, err := price.New(&stripe.PriceParams{
		Currency: stripe.String("usd"),
		Product:  stripe.String(subscriptionProd.ID),
		Recurring: &stripe.PriceRecurringParams{
			Interval: stripe.String("month"),
		},
		UnitAmount: stripe.Int64(request.Price),
	})

	if err != nil {
		fmt.Printf("\nError when creating subscription price for product on stripe: %+v\n\n", err)
		return nil, err

	}

	fmt.Printf("Created new subscription's price successfully. Response:%+v\n", subscriptionPrice)

	// set default price. NOT set by default.
	product.Update(subscriptionProd.ID, &stripe.ProductParams{
		DefaultPrice: stripe.String(subscriptionPrice.ID),
	})

	return &SetupProductsResp{
		PriceID: subscriptionPrice.ID,
	}, nil

}

/**
* Lists all the products existing on the stripe catalog from all customers.
**/
func (s *StripeProcessor) GetProducts(ctx context.Context) (*ProductListResponse, error) {
	// gets all active products with prices
	params := &stripe.ProductListParams{
		Active: stripe.Bool(true),
	}
	params.AddExpand("data.default_price")

	iter := product.List(params)

	var productList []ProductInfo
	for iter.Next() {
		prod := iter.Product()

		// skip products without a default price
		if prod.DefaultPrice == nil {
			continue
		}

		productInfo := ProductInfo{
			ID:          prod.ID,
			Name:        prod.Name,
			Description: prod.Description,
		}

		// Get price information from the expanded default_price
		if prod.DefaultPrice != nil {
			productInfo.PriceID = prod.DefaultPrice.ID
			productInfo.Price = prod.DefaultPrice.UnitAmount

			// Determine if it's a subscription or one-time product
			if prod.DefaultPrice.Recurring != nil {
				productInfo.Type = "subscription"
			} else {
				productInfo.Type = "one-time"
			}
		}

		productList = append(productList, productInfo)
	}

	// fmt.Printf("\nproductList: %+v\n\n", productList)

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("error listing products: %w", err)
	}

	return &ProductListResponse{Products: productList}, nil
}

/**
* Purchases a specific product by creating a payment intent for the product's price.
**/
func (s *StripeProcessor) PurchaseProduct(ctx context.Context, req *PurchaseProductRequest) (*StripePurchaseResponse, error) {
	// first, get the product to find its default price
	productParams := &stripe.ProductParams{}

	// we need to use AddExpand method on the the field we want to convert
	// from an id to the object with DETAILED data object.
	//
	// Normal Stripe Response:
	// json{
	//   "id": "prod_123",
	//   "name": "T-Shirt",
	//   "default_price": "price_456"  // Just a string ID
	// }
	//
	// With AddExpand("default_price"):
	// json{
	//   "id": "prod_123",
	//   "name": "T-Shirt",
	//   "default_price": {  // Full object!
	//     "id": "price_456",
	//     "unit_amount": 2000,
	//     "currency": "usd",
	//     "recurring": null
	//   }
	// }

	productParams.AddExpand("default_price")

	prod, err := product.Get(req.ProductID, productParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if prod.DefaultPrice == nil {
		return nil, fmt.Errorf("product has no default price")
	}

	fmt.Printf("Product price amount: %d\n", prod.DefaultPrice.UnitAmount)

	// create payment intent with the product's price
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(prod.DefaultPrice.UnitAmount),
		Currency: stripe.String("usd"),
		Customer: stripe.String(req.CustomerID),

		ConfirmationMethod: stripe.String("automatic"),
		Confirm:            stripe.Bool(false),

		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),

		// add metadata to track the product being purchased
		Metadata: map[string]string{
			"product_id": req.ProductID,
			"price_id":   prod.DefaultPrice.ID,
		},
	}

	intent, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	return &StripePurchaseResponse{
		ClientSecret:    intent.ClientSecret,
		PaymentIntentID: intent.ID,
		Amount:          prod.DefaultPrice.UnitAmount,
	}, nil
}

/**
* Subscribes a specific product by creating a payment intent for the product's price.
**/
func (s *StripeProcessor) SubscribeToProduct(ctx context.Context, req *SubscribeRequest) (*SubscribeResponse, error) {
	// get product with expanded default_price to check if it's a subscription
	productParams := &stripe.ProductParams{}
	productParams.AddExpand("default_price")

	prod, err := product.Get(req.ProductID, productParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if prod.DefaultPrice == nil {
		return nil, fmt.Errorf("product has no default price")
	}

	if prod.DefaultPrice.Recurring == nil {
		return nil, fmt.Errorf("product %s is not a subscription (no recurring price)", req.ProductID)
	}

	// create subscription
	subParams := &stripe.SubscriptionParams{
		Customer: stripe.String(req.CustomerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(prod.DefaultPrice.ID),
			},
		},
		PaymentBehavior: stripe.String("default_incomplete"),
		PaymentSettings: &stripe.SubscriptionPaymentSettingsParams{
			SaveDefaultPaymentMethod: stripe.String("on_subscription"),
		},
	}

	subParams.AddExpand("latest_invoice.confirmation_secret")

	// create the subscription
	sub, err := subscription.New(subParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	// extract client secret from the invoice's confirmation_secret
	var clientSecret string
	if sub.LatestInvoice != nil && sub.LatestInvoice.ConfirmationSecret != nil {
		// the ConfirmationSecret contains the client_secret
		clientSecret = sub.LatestInvoice.ConfirmationSecret.ClientSecret
	}

	return &SubscribeResponse{
		SubscriptionID: sub.ID,
		ClientSecret:   clientSecret,
		Status:         string(sub.Status),
	}, nil
}

/**
* Checks webhook event, using it as an indicator that something has been triggered.
* if you are on my team and I provide you this POC as guidance please be wary that the actual "sync" method comes after this processor
**/

func (s *StripeProcessor) ProcessWebhookEvent(ctx context.Context, event *stripe.Event) (customerId string, error error) {
	isEventSupported := s.IsWebhookEventSupported(ctx, event)

	if !isEventSupported {
		return "", fmt.Errorf("The event type that was resulted from the action was not allowed.")
	}

	customerId, err := s.ExtractCustomerIdFromWebhook(event)
	fmt.Printf("customerId from webhook event: %s\n", customerId)

	if err != nil {
		return "", err
	}

	return customerId, nil
}

/**
* Handles all webhook events specifically for stripe that this platform allows.
**/

func (s *StripeProcessor) IsWebhookEventSupported(ctx context.Context, event *stripe.Event) bool {
	// store allowed / expected webhook events
	expectedEvents := map[stripe.EventType]bool{
		stripe.EventTypePaymentIntentSucceeded:      true,
		stripe.EventTypePaymentIntentPaymentFailed:  true,
		stripe.EventTypePaymentIntentCanceled:       true,
		stripe.EventTypeCustomerSubscriptionCreated: true,
	}

	fmt.Printf("Processing webhook event type: %s\n", event.Type)

	if !expectedEvents[event.Type] {
		// not-allowed events
		fmt.Printf("The event type that was resulted from the action was not allowed.")
		return false
	}

	return true
}

/**
* Extract stripe-specific customer id from event object.
**/
func (s *StripeProcessor) ExtractCustomerIdFromWebhook(event interface{}) (string, error) {
	// Type assert to Stripe event
	stripeEvent, ok := event.(*stripe.Event)
	if !ok {
		fmt.Printf("invalid event type: expected *stripe.Event")
		return "", fmt.Errorf("invalid event type: expected *stripe.Event")
	}

	var eventData map[string]interface{}
	err := json.Unmarshal(stripeEvent.Data.Raw, &eventData)

	if err != nil {
		fmt.Printf("\nCould not unmarshal event data into eventData, err: %+v\n\n", err)
		return "", err
	}

	if customer, ok := eventData["customer"].(string); ok && customer != "" {
		fmt.Printf("\ncustomer from event: %s\n\n", customer)
		return customer, nil
	}

	return "", fmt.Errorf("no customer ID found in stripe event type: %s", stripeEvent.Type)
}

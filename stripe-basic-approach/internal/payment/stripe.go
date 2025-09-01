package payment

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82/price"
	"github.com/stripe/stripe-go/v82/product"
	"github.com/stripe/stripe-go/v82/setupintent"
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

	// create SUBSCRIPTION product / service
	// subscriptionPrice, err := price.New(&stripe.PriceParams{
	// 	Currency: stripe.String("usd"),
	// 	Product:  stripe.String(product.ID),
	// 	Recurring: &stripe.PriceRecurringParams{
	// 		Interval: stripe.String("month"),
	// 	},
	// 	UnitAmount: stripe.Int64(request.Price),
	// })

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

func (s *StripeProcessor) CreateSubscription(ctx context.Context, request *SetupProductsReq) (*SetupProductsResp, error) {
	// create product
	prod, err := product.New(&stripe.ProductParams{
		Name:        stripe.String(request.Name),
		Description: stripe.String(request.Description),
	})

	if err != nil {
		fmt.Printf("\nError when creating product on stripe: %+v\n\n", err)
		return nil, err
	}

	// create SUBSCRIPTION product / service
	// subscriptionPrice, err := price.New(&stripe.PriceParams{
	// 	Currency: stripe.String("usd"),
	// 	Product:  stripe.String(product.ID),
	// 	Recurring: &stripe.PriceRecurringParams{
	// 		Interval: stripe.String("month"),
	// 	},
	// 	UnitAmount: stripe.Int64(request.Price),
	// })

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
		SubscriptionPriceID: oneTimePrice.ID,
	}, nil
}
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
		}

		productList = append(productList, productInfo)
	}

	fmt.Printf("\nproductList: %+v\n\n", productList)

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("error listing products: %w", err)
	}

	return &ProductListResponse{Products: productList}, nil
}

/**
* Purchases a specific product by creating a payment intent for the product's price.
**/
func (s *StripeProcessor) PurchaseProduct(ctx context.Context, req *PurchaseProductRequest) (*PurchaseProductResponse, error) {
	// first, get the product to find its default price
	prod, err := product.Get(req.ProductID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	if prod.DefaultPrice == nil {
		return nil, fmt.Errorf("product has no default price")
	}

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

	return &PurchaseProductResponse{
		ClientSecret:    intent.ClientSecret,
		PaymentIntentID: intent.ID,
	}, nil
}

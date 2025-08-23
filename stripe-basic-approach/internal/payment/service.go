package payment

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/price"
	"github.com/stripe/stripe-go/v82/product"
	"github.com/stripe/stripe-go/v82/setupintent"
)

type service struct {
	userService PaymentUserService
}

type PaymentUserService interface {
	UpdateStripeCustomer(ctx context.Context, userID uuid.UUID, stripeCustomerID string) error
}

func NewService() *service {
	return &service{}
}

/**
* Sets up custom products for purchase. This can be anything you wish to sell, and is a digital representation of the
* item - which can be physical, digital, or even just a concept (donation, etc).
**/
func (s *service) SetupProducts(ctx context.Context, request *SetupProductsReq) (*SetupProductsResp, error) {
	// create product
	oneTimeProduct, err := product.New(&stripe.ProductParams{
		Name:        stripe.String(request.Name),
		Description: stripe.String(request.Description),
	})

	if err != nil {
		fmt.Printf("\nError when creating product on stripe: %+v\n\n", err)
		return nil, err
	}

	// create subscription price (recurring monthly)
	subscriptionPrice, err := price.New(&stripe.PriceParams{
		Currency: stripe.String("usd"),
		Product:  stripe.String(oneTimeProduct.ID),
		Recurring: &stripe.PriceRecurringParams{
			Interval: stripe.String("month"),
		},
		UnitAmount: stripe.Int64(request.Price),
	})

	if err != nil {
		fmt.Printf("\nError when creating subscription price for product on stripe: %+v\n\n", err)
		return nil, err
	}

	fmt.Printf("Created new product's price successfully. Response:%+v\n", subscriptionPrice)

	return &SetupProductsResp{
		SubscriptionPriceID: subscriptionPrice.ID,
	}, nil
}

func (s *service) CreateCustomer(ctx context.Context, userId uuid.UUID, email string) (string, error) {
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

	log.Printf("Created customer: %s", cust.ID)
	return cust.ID, nil
}

/**
* This method AUTHORIZES a card save for a customer by creating a permission token for the client
* to then use the stripe sdk via elements to save the card.
**/
func (s *service) SaveCard(ctx context.Context, customerId string) (string, error) {
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

func (s *service) CreatePaymentIntent(ctx context.Context, amount int64, customerId string) (string, error) {
	return "", nil
}

func (s *service) CreateSubscription(ctx context.Context, priceId, customerId, email string) (*SubscriptionResp, error) {
	return nil, nil
}

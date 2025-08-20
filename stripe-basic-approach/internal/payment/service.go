package payment

import (
	"context"
	"fmt"
	"log"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/price"
	"github.com/stripe/stripe-go/v82/product"
)

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) SetupProducts(ctx context.Context, request *SetupProductsReq) error {
	// create product
	oneTimeProduct, err := product.New(&stripe.ProductParams{
		Name:        stripe.String(request.Name),
		Description: stripe.String(request.Description),
	})

	if err != nil {
		return err
	}

	// create one-time price
	oneTimePrice, err := price.New(&stripe.PriceParams{
		Currency:   stripe.String("usd"),
		Product:    stripe.String(oneTimeProduct.ID),
		UnitAmount: stripe.Int64(request.Price),
	})

	if err != nil {
		return err
	}

	fmt.Printf("Created new product's price successfully. Response:%+v\n", oneTimePrice)

	return nil
}

func (s *service) CreateCustomer(context.Context, *CreateCustomerReq) error {
	params := &stripe.CustomerParams{
		Email: stripe.String("demo@example.com"),
		Metadata: map[string]string{
			"user_id": "demo_user_123", // Link to your user system
		},
	}

	// create customer on Stripe
	cust, err := customer.New(params)
	if err != nil {
		return err
	}

	log.Printf("Created customer: %s", cust.ID)
	return nil
}

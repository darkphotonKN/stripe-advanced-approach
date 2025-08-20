package payment

import (
	"context"
	"log"

	"github.com/docker/distribution/uuid"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
)

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) SetupProducts(userId uuid.UUID) error {
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

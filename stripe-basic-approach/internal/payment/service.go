package payment

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type service struct {
	userService      PaymentUserService
	paymentProcessor PaymentProcessor
	repo             Repository
}

// SOLI"D"
// DIP

type Repository interface {
	Create(ctx context.Context, request *CheckoutSessionRequest) (uuid.UUID, error)
}

type PaymentUserService interface {
	UpdateStripeCustomer(ctx context.Context, userID uuid.UUID, stripeCustomerID string) error
}

func NewService(userService PaymentUserService, paymentProcessor PaymentProcessor) *service {
	return &service{
		userService:      userService,
		paymentProcessor: paymentProcessor,
	}
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

func (s *service) PurchaseProduct(ctx context.Context, req *PurchaseProductRequest) (*PurchaseProductResponse, error) {
	return s.paymentProcessor.PurchaseProduct(ctx, req)
}

func (s *service) SetupSubscription(ctx context.Context, request *SetupProductsReq) (*SetupProductsResp, error) {
	return s.paymentProcessor.SetupSubscription(ctx, request)
}

func (s *service) SubscribeToProduct(ctx context.Context, req *SubscribeRequest) (*SubscribeResponse, error) {
	return s.paymentProcessor.SubscribeToProduct(ctx, req)
}

/*
	--- Full checkout session flow, for handling Payment Success / Failure ---

	User Journey:
	- User completes payment/subscription
	- Redirects back to your app's /success page
	- Success page shows "Payment processing..."
	- Waits for webhook to update database
	- Eventually shows subscription status
*/

/**
* Endpoint 1 - Creates Stripe session, saves pending payment
**/

func (s *service) CreateCheckoutSession(ctx context.Context, req *CheckoutSessionRequest) (*CheckoutSessionResponse, error) {
	// 1. Create pending payment record in database
	paymentID, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment record: %w", err)
	}

	// 2. Create Stripe checkout session
	session, err := s.paymentProcessor.CreateCheckoutSession(ctx, req.CustomerID, paymentID)
	if err != nil {
		return nil, err
	}

	// 3. Update payment record with session ID
	// _, err = s.db.ExecContext(ctx, `
	//       UPDATE payments
	//       SET stripe_session_id = $1
	//       WHERE id = $2
	//   `, session.SessionID, paymentID)
	//
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to update payment: %w", err)
	// }

	return &CheckoutSessionResponse{
		SessionID:   session.SessionID,
		CheckoutURL: session.CheckoutURL,
	}, nil
}

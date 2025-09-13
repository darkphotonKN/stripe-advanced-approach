package payment

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
)

type service struct {
	userService      PaymentUserService
	paymentProcessor PaymentProcessor
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

func NewService(repo Repository, userService PaymentUserService, paymentProcessor PaymentProcessor) *service {
	return &service{
		repo:             repo,
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

func (s *service) SubscribeToProduct(ctx context.Context, req *SubscribeRequest) (*SubscribeResponse, error) {
	return s.paymentProcessor.SubscribeToProduct(ctx, req)
}

// --- Full Flow Methods ---

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

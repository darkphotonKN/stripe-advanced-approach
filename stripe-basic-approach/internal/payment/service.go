package payment

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type service struct {
	userService      PaymentUserService
	paymentProcessor PaymentProcessor
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

func (s *service) CreateSubscription(ctx context.Context, priceId, customerId, email string) (*SubscriptionResp, error) {
	return nil, nil
}

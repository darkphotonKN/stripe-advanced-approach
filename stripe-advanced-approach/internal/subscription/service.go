package subscription

import (
	"context"

	"github.com/darkphotonKN/stripe-advanced-approach/internal/model"
	"github.com/darkphotonKN/stripe-advanced-approach/internal/user"
	"github.com/google/uuid"
)

type service struct {
	userService    SubscriptionUserService
	paymentService SubscriptionPaymentService
}

type SubscriptionUserService interface {
	Update(ctx context.Context, id uuid.UUID, user *user.User) error
}

type SubscriptionPaymentService interface {
	GetSubscriptionStatusCache(ctx context.Context, userId uuid.UUID) (*model.SubscriptionStatus, error)
}

func NewService(userService SubscriptionUserService, paymentService SubscriptionPaymentService) Service {
	return &service{
		userService:    userService,
		paymentService: paymentService,
	}
}

// SubscribeToProduct creates a subscription for a user
func (s *service) SubscribeToProduct(ctx context.Context, userId uuid.UUID, req *SubscribeRequest) (*SubscribeResponse, error) {
	return nil, nil
}

// GetSubscriptionStatus retrieves the user's subscription status
func (s *service) GetSubscriptionStatus(ctx context.Context, userId uuid.UUID) (*SubscriptionStatusResponse, error) {
	// TODO: complete implementation
	_, err := s.paymentService.GetSubscriptionStatusCache(ctx, userId)
	return nil, err
}

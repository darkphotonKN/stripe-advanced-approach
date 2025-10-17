package subscription

import (
	"context"

	"github.com/google/uuid"
)

type service struct {
	repo Repository
}

type Repository interface {
	// TODO: Add repository methods as needed
	// Example: GetSubscriptionByUserID(ctx context.Context, userId uuid.UUID) (*Subscription, error)
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// SubscribeToProduct creates a subscription for a user
// TODO: Implement this method
func (s *service) SubscribeToProduct(ctx context.Context, userId uuid.UUID, req *SubscribeRequest) (*SubscribeResponse, error) {
	// TODO: Implement subscription logic
	// 1. Get user's Stripe customer ID
	// 2. Create Stripe subscription or checkout session
	// 3. Return client secret or checkout URL
	return nil, nil
}

// GetSubscriptionStatus retrieves the user's subscription status
// TODO: Implement this method
func (s *service) GetSubscriptionStatus(ctx context.Context, userId uuid.UUID) (*SubscriptionStatusResponse, error) {
	// TODO: Implement status check logic
	// 1. Get user's Stripe customer ID
	// 2. Fetch subscription from Stripe
	// 3. Map to minimal frontend response
	return nil, nil
}

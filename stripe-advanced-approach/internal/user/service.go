package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByStripeCustomerID(ctx context.Context, stripeCustomerID string) (*User, error)
	List(ctx context.Context) ([]User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStripeCustomer(ctx context.Context, userID uuid.UUID, stripeCustomerID string) error
}

type service struct {
	repo           Repository
	paymentService UserPaymentService
}

type UserPaymentService interface {
	AddCacheUserIdToCusId(ctx context.Context, userId uuid.UUID, customerId string) error
	AddCacheCusIdToUserId(ctx context.Context, customerId string, userId uuid.UUID) error
	CreateCustomer(ctx context.Context, userId uuid.UUID, email string) (string, error)
	SyncStripeDataToStorage(ctx context.Context, customerId string) error
}

func NewService(repo Repository) *service {
	return &service{
		repo: repo,
	}
}

/**
* dependency injection for payment service after payment service finished setting up with base user service
**/
func (s *service) SetPaymentService(paymentService UserPaymentService) {
	s.paymentService = paymentService
}

func (s *service) Create(ctx context.Context, user *User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	createdUser, err := s.repo.Create(ctx, user)

	if err != nil {
		fmt.Printf("could not create user, err:%s\n", err)
		return err
	}

	// create a payment processor user once user is created on platform
	customerID, err := s.paymentService.CreateCustomer(ctx, createdUser.ID, user.Email)

	if err != nil {
		fmt.Printf("could not create payment processor customer.\n")
		return err
	}

	// sync to cache
	go s.SyncCacheAndMappings(ctx, createdUser.ID, customerID)
	return nil

}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid ID")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetByEmail(ctx context.Context, email string) (*User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	return s.repo.GetByEmail(ctx, email)
}

func (s *service) List(ctx context.Context) ([]User, error) {
	return s.repo.List(ctx)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, user *User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	user.ID = id
	return s.repo.Update(ctx, user)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid ID")
	}
	return s.repo.Delete(ctx, id)
}

func (s *service) Authenticate(ctx context.Context, email, password string) (*User, error) {
	// --- User Authentication ---
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		fmt.Printf("Error when authenticating email: %s\n", err)
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	user.Password = ""

	// --- Cache Updates ---

	customerId := *user.StripeCustomerID
	userId := user.ID

	// don't stop the flow, but log the failure of cache updates
	// BUG: still not working
	// TODO: fix
	go s.SyncCacheAndMappings(ctx, userId, customerId)

	return user, nil
}

func (s *service) SyncCacheAndMappings(ctx context.Context, userId uuid.UUID, customerId string) {
	// updates cache with customerId to userId mapping
	err := s.paymentService.AddCacheUserIdToCusId(ctx, userId, customerId)

	if err != nil {
		fmt.Printf("Error syncing up userId to customerId cache: %s\n", err)
	}

	// updates cache with customerId to userId mapping
	err = s.paymentService.AddCacheCusIdToUserId(ctx, customerId, userId)

	if err != nil {
		fmt.Printf("Error syncing up customerId to userId in cache: %s\n", err)
	}

	// sync primary customer data in cache
	err = s.paymentService.SyncStripeDataToStorage(ctx, customerId)

	if err != nil {
		fmt.Printf("Error syncing up customer data in cache: %s\n", err)
	}
}

func (s *service) UpdateStripeCustomer(ctx context.Context, userId uuid.UUID, customerId string) error {
	// not ok to fail
	// don't update cache until repo is updated successfully to prevent unnecessary rollbacks
	err := s.repo.UpdateStripeCustomer(ctx, userId, customerId)

	if err != nil {
		return err
	}

	// updates cache with customerId to userId mapping
	go s.SyncCacheAndMappings(ctx, userId, customerId)

	return nil
}

func (s *service) GetStripeCustomer(ctx context.Context, userID uuid.UUID) (*string, error) {
	if userID == uuid.Nil {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user.StripeCustomerID, nil
}

func (s *service) GetByStripeCustomerID(ctx context.Context, stripeCustomerID string) (*User, error) {
	if stripeCustomerID == "" {
		return nil, errors.New("stripe customer ID is required")
	}
	return s.repo.GetByStripeCustomerID(ctx, stripeCustomerID)
}

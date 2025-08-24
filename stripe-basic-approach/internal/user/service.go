package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context) ([]User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStripeCustomer(ctx context.Context, userID uuid.UUID, stripeCustomerID string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
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

	return s.repo.Create(ctx, user)
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
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		fmt.Printf("Error when authenticating email: %s\n", err)
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	user.Password = ""
	return user, nil
}

func (s *service) UpdateStripeCustomer(ctx context.Context, userId uuid.UUID, customerId string) error {

	return s.repo.UpdateStripeCustomer(ctx, userId, customerId)

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

package product

import (
	"context"
	"errors"
)

type Service interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id int) (*Product, error)
	List(ctx context.Context) ([]Product, error)
	Update(ctx context.Context, id int, product *Product) error
	Delete(ctx context.Context, id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, product *Product) error {
	if err := product.Validate(); err != nil {
		return err
	}
	return s.repo.Create(ctx, product)
}

func (s *service) GetByID(ctx context.Context, id int) (*Product, error) {
	if id <= 0 {
		return nil, errors.New("invalid ID")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *service) List(ctx context.Context) ([]Product, error) {
	return s.repo.List(ctx)
}

func (s *service) Update(ctx context.Context, id int, product *Product) error {
	if err := product.Validate(); err != nil {
		return err
	}
	product.ID = id
	return s.repo.Update(ctx, product)
}

func (s *service) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid ID")
	}
	return s.repo.Delete(ctx, id)
}


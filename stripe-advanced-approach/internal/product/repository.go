package product

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id int) (*Product, error)
	List(ctx context.Context) ([]Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id int) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, product *Product) error {

	query := `
		INSERT INTO products (name, description, price, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.db.GetContext(ctx, product, query, product.Name, product.Description, product.Price, product.Stock)
}

func (r *repository) GetByID(ctx context.Context, id int) (*Product, error) {
	var product Product
	query := `SELECT * FROM products WHERE id = $1`
	err := r.db.GetContext(ctx, &product, query, id)
	return &product, err
}

func (r *repository) List(ctx context.Context) ([]Product, error) {
	var products []Product
	query := `SELECT * FROM products ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &products, query)
	return products, err
}

func (r *repository) Update(ctx context.Context, product *Product) error {
	query := `
		UPDATE products
		SET name = $1, description = $2, price = $3, stock = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING updated_at
	`
	return r.db.GetContext(ctx, product, query, product.Name, product.Description, product.Price, product.Stock, product.ID)
}

func (r *repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}


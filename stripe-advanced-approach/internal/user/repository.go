package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user *User) (*User, error) {
	fmt.Printf("repo create inc user: %+v\n", *user)
	fmt.Printf("repo create inc user.email %s\n user.password %s \nuser.name: %+v\n", user.Email, user.Password, user.Name)

	query := `
		INSERT INTO users (email, password, name, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, email, name, created_at, updated_at
	`
	var createdUser User
	err := r.db.GetContext(ctx, &createdUser, query, user.Email, user.Password, user.Name)

	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.db.GetContext(ctx, &user, query, id)
	return &user, err
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query := `SELECT * FROM users WHERE email = $1`
	err := r.db.GetContext(ctx, &user, query, email)
	return &user, err
}

func (r *repository) List(ctx context.Context) ([]User, error) {
	var users []User
	query := `SELECT * FROM users ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &users, query)
	return users, err
}

func (r *repository) Update(ctx context.Context, user *User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, subscribed = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at
	`
	return r.db.GetContext(ctx, user, query, user.Name, user.Email, user.Subscribed, user.ID)
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *repository) UpdateStripeCustomer(ctx context.Context, userID uuid.UUID, stripeCustomerID string) error {
	query := `
		UPDATE users
		SET stripe_customer_id = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, stripeCustomerID, userID)
	return err
}

func (r *repository) GetByStripeCustomerID(ctx context.Context, stripeCustomerID string) (*User, error) {
	var user User
	query := `SELECT * FROM users WHERE stripe_customer_id = $1`
	err := r.db.GetContext(ctx, &user, query, stripeCustomerID)
	return &user, err
}

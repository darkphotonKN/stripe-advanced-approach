package user

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID               uuid.UUID `db:"id" json:"id"`
	Email            string    `db:"email" json:"email" validate:"required,email"`
	Password         string    `db:"password" json:"password,omitempty" validate:"required,min=6"`
	Name             string    `db:"name" json:"name" validate:"required,min=1,max=255"`
	StripeCustomerID *string   `db:"stripe_customer_id" json:"stripe_customer_id"`
	Subscribed       bool      `db:"subscribed" json:"subscribed"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

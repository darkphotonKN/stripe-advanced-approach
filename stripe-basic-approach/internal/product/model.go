package product

import (
	"time"
	"github.com/go-playground/validator/v10"
)

type Product struct {
	ID          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name" validate:"required,min=1,max=255"`
	Description string    `db:"description" json:"description"`
	Price       float64   `db:"price" json:"price" validate:"required,min=0"`
	Stock       int       `db:"stock" json:"stock" validate:"min=0"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

func (p *Product) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
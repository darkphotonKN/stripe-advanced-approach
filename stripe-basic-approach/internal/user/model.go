package user

import (
	"time"
	"github.com/go-playground/validator/v10"
)

type User struct {
	ID        int       `db:"id" json:"id"`
	Email     string    `db:"email" json:"email" validate:"required,email"`
	Password  string    `db:"password" json:"password,omitempty" validate:"required,min=6"`
	Name      string    `db:"name" json:"name" validate:"required,min=1,max=255"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
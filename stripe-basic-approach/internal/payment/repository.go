package payment

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, temp *CheckoutSessionRequest) (uuid.UUID, error) {
	// _, err := r.db.ExecContext(ctx, `
	//       INSERT INTO payments (
	//           stripe_customer_id,
	//           status,
	//           created_at,
	//           updated_at
	//       ) VALUES ($1, $2, NOW(), NOW())
	//   `, req.paymentID, req.CustomerID, "pending")
	//
	return uuid.Nil, nil
}

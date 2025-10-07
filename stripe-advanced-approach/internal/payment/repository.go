package payment

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

func (r *repository) Create(ctx context.Context, userId uuid.UUID, paymentIntent *PaymentIntentRequest) error {

	query := `
		INSERT INTO payments (
			user_id,
			stripe_customer_id,
			stripe_payment_intent_id,
			amount, 
			status,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	`

	_, err := r.db.ExecContext(ctx, query, userId, paymentIntent.CustomerID, paymentIntent.IntentID, paymentIntent.Amount, "pending")

	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetPaymentByIntentID(ctx context.Context, intentID string) (*Payment, error) {
	var payment Payment

	query := `
		SELECT * FROM payments
		WHERE stripe_intent_id = $1
	`

	err := r.db.GetContext(ctx, &payment, query, intentID)
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *repository) UpdateStatus(ctx context.Context, intentID string, status string) error {
	query := `
		UPDATE payments
		SET status = $1, updated_at = NOW()
		WHERE stripe_payment_intent_id = $2
	`

	_, err := r.db.ExecContext(ctx, query, status, intentID)

	if err != nil {
		fmt.Printf("\nError when updating payment table status column: %+v\n\n", err)
		return err
	}

	return nil
}

func (r *repository) UpsertPayment(ctx context.Context, paymentIntentID string, payment *Payment) error {
	query := `
        INSERT INTO payments (
            user_id,
            stripe_customer_id,
            stripe_payment_intent_id,
            amount,
            status,
            currency,
            created_at,
            updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
        ON CONFLICT (stripe_payment_intent_id) 
        DO UPDATE SET
            amount = EXCLUDED.amount,
            status = EXCLUDED.status,
            currency = EXCLUDED.currency,
            updated_at = NOW()
        RETURNING id
    `

	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, query,
		payment.UserID,
		payment.StripeCustomerID,
		paymentIntentID,
		payment.Amount,
		payment.Status,
		payment.Currency,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("failed to upsert payment: %w", err)
	}

	return nil
}

func (r *repository) UpsertSubscriptionRecord(ctx context.Context, sub *Subscription) error {
	query := `
		INSERT INTO subscriptions (
			user_id,
			stripe_customer_id,
			stripe_subscription_id,
			stripe_price_id,
			status,
			current_period_start,
			current_period_end,
			cancel_at_period_end,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		ON CONFLICT (stripe_subscription_id)
		DO UPDATE SET
			stripe_price_id = EXCLUDED.stripe_price_id,
			status = EXCLUDED.status,
			current_period_start = EXCLUDED.current_period_start,
			current_period_end = EXCLUDED.current_period_end,
			cancel_at_period_end = EXCLUDED.cancel_at_period_end,
			updated_at = NOW()
		RETURNING id
	`

	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, query,
		sub.UserID,
		sub.StripeCustomerID,
		sub.StripeSubscriptionID,
		sub.StripePriceID,
		sub.Status,
		sub.CurrentPeriodStart,
		sub.CurrentPeriodEnd,
		sub.CancelAtPeriodEnd,
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("failed to upsert subscription: %w", err)
	}

	return nil
}

func (r *repository) GetActiveSubscription(ctx context.Context, userID uuid.UUID) (*Subscription, error) {
	var subscription Subscription

	query := `
		SELECT * FROM subscriptions
		WHERE user_id = $1 AND status = 'active'
		ORDER BY created_at DESC
		LIMIT 1
	`

	err := r.db.GetContext(ctx, &subscription, query, userID)
	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (r *repository) UpdateSubscriptionStatus(ctx context.Context, subID string, status string) error {
	query := `
		UPDATE subscriptions
		SET status = $1, updated_at = NOW()
		WHERE stripe_subscription_id = $2
	`

	_, err := r.db.ExecContext(ctx, query, status, subID)

	if err != nil {
		fmt.Printf("\nError when updating subscription status %+v\n\n", err)
		return err
	}

	return nil
}

func (r *repository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}

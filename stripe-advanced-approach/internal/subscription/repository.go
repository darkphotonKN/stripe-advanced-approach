package subscription

import (
	"github.com/jmoiron/sqlx"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

// TODO: Implement repository methods as needed
// Examples:
// - GetSubscriptionByUserID
// - CreateSubscription
// - UpdateSubscription
// - DeleteSubscription

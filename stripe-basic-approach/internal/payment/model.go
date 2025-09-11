package payment

import (
	"time"

	"github.com/google/uuid"
)

// General Payments Entity
type Payment struct {
	ID                 uuid.UUID  `db:"id" json:"id"`
	UserID             uuid.UUID  `db:"user_d" json:"user_id"`
	StripeCustomerID   string     `db:"stripe_customer_id" json:"stripe_customer_id"`
	StripeIntentID     string     `db:"stripe_intent_id" json:"stripe_intent_id"`
	StripeSessionID    string     `db:"stripe_session_id" json:"stripe_session_id"`
	Amount             int64      `db:"amount" json:"amount"`
	Currency           string     `db:"currency" json:"currency"`
	Status             string     `db:"status" json:"status"`
	PaymentMethodTypes string     `db:"payment_method_types" json:"payment_method_types"`
	CreatedAt          time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at" json:"updated_at"`
	CompletedAt        *time.Time `db:"completed_at" json:"completed_at"`
}

// Subscription Entity
type Subscription struct {
	ID                   uuid.UUID  `db:"id" json:"id"`
	UserID               uuid.UUID  `db:"user_id" json:"user_id"`
	StripeCustomerID     string     `db:"stripe_customer_id" json:"stripe_customer_id"`
	StripeSubscriptionID string     `db:"stripe_subscription_id" json:"stripe_subscription_id"`
	StripePriceID        string     `db:"stripe_price_id" json:"stripe_price_id"`
	Status               string     `db:"status" json:"status"`
	CurrentPeriodStart   time.Time  `db:"current_period_start" json:"current_period_start"`
	CurrentPeriodEnd     time.Time  `db:"current_period_end" json:"current_period_end"`
	CancelAtPeriodEnd    bool       `db:"cancel_at_period_end" json:"cancel_at_period_end"`
	CanceledAt           *time.Time `db:"canceled_at" json:"canceled_at"`
	CreatedAt            time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time  `db:"updated_at" json:"updated_at"`
}

// Setup Products
type SetupProductsReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
}

type SetupProductsResp struct {
	PriceID string `json:"price_id"`
}

// Create Customer

type CreateCustomerReq struct {
	Email string `json:"email"`
}

type CreateCustomerResponse struct {
	CustomerID string `json:"customerId"`
}

type CreatePaymentIntentRequest struct {
	Amount     int64  `json:"amount"`
	CustomerID string `json:"customer_id"`
}

type CreatePaymentIntentResponse struct {
	ClientSecret    string `json:"client_secret"`
	PaymentIntentID string `json:"paymentIntentId"`
}

type SaveCardRequest struct {
	CustomerID string `json:"customerId"`
}

type SaveCardResponse struct {
	ClientSecret  string `json:"client_secret"`
	SetupIntentID string `json:"setupIntentId"`
}

type SubscriptionResp struct {
	SubscriptionID string `json:"subscription_id"`
	ClientSecret   string `json:"client_secret"`
}

// Products List
type ProductInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	PriceID     string `json:"price_id"`
	Type        string `json:"type"` // "one-time" or "subscription"
}

type ProductListResponse struct {
	Products []ProductInfo `json:"products"`
}

// Purchase Product
type PurchaseProductRequest struct {
	ProductID  string `json:"product_id" binding:"required"`
	CustomerID string `json:"customer_id" binding:"required"`
}

type PurchaseProductResponse struct {
	ClientSecret    string `json:"client_secret"`
	PaymentIntentID string `json:"payment_intent_id"`
}

// Internal Stripe response type with additional priceID
type StripePurchaseResponse struct {
	ClientSecret    string `json:"client_secret"`
	PaymentIntentID string `json:"payment_intent_id"`
	Amount          int64  `json:"amount"`
}

// Subscribe Product
type SubscribeRequest struct {
	ProductID  string `json:"product_id"`  // Product to subscribe to
	CustomerID string `json:"customer_id"` // Stripe customer ID
}

type SubscribeResponse struct {
	SubscriptionID string `json:"subscription_id"` // sub_xxx ID for management
	ClientSecret   string `json:"client_secret"`   // For frontend to confirm payment
	Status         string `json:"status"`          // "incomplete" until payment confirmed
}

// Payment Intent Request for internal use
type PaymentIntentRequest struct {
	CustomerID string `json:"customer_id" db:"customer_id"`
	Amount     int64  `json:"amount" db:"amount"`
	IntentID   string `json:"intent_id" db:"stripe_intent_id"`
}

type CheckoutSessionResponse struct {
	SessionID   string `json:"session_id"`
	CheckoutURL string `json:"checkout_url"`
}

// Success
type SuccessReponse struct {
}

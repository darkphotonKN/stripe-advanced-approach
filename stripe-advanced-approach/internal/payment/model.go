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
	Status             string     `db:"status" json:"status"` // synced from stripe
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

// stripe customer cached data
type StripeCacheData struct {
	ID                   string                   `json:"id"`
	Address              *CustomerAddress         `json:"address"`
	Balance              int64                    `json:"balance"`
	CashBalance          *CustomerCashBalance     `json:"cash_balance"`
	Created              int64                    `json:"created"`
	Currency             string                   `json:"currency"`
	DefaultSource        *string                  `json:"default_source"`
	Deleted              bool                     `json:"deleted"`
	Delinquent           bool                     `json:"delinquent"`
	Description          string                   `json:"description"`
	Discount             *CustomerDiscount        `json:"discount"`
	Email                string                   `json:"email"`
	InvoiceCreditBalance map[string]int64         `json:"invoice_credit_balance"`
	InvoicePrefix        string                   `json:"invoice_prefix"`
	InvoiceSettings      *CustomerInvoiceSettings `json:"invoice_settings"`
	Livemode             bool                     `json:"livemode"`
	Metadata             map[string]string        `json:"metadata"`
	Name                 string                   `json:"name"`
	NextInvoiceSequence  int64                    `json:"next_invoice_sequence"`
	Object               string                   `json:"object"`
	Phone                string                   `json:"phone"`
	PreferredLocales     []string                 `json:"preferred_locales"`
	Shipping             *CustomerShipping        `json:"shipping"`
	Sources              *CustomerSourceList      `json:"sources"`
	Subscriptions        *SubscriptionList        `json:"subscriptions"`
	Tax                  *CustomerTax             `json:"tax"`
	TaxExempt            string                   `json:"tax_exempt"`
	TaxIds               *CustomerTaxIdList       `json:"tax_ids"`
	TestClock            *TestClock               `json:"test_clock"`
}

type CustomerAddress struct {
	City       string `json:"city"`
	Country    string `json:"country"`
	Line1      string `json:"line1"`
	Line2      string `json:"line2"`
	PostalCode string `json:"postal_code"`
	State      string `json:"state"`
}

type CustomerCashBalance struct {
	Object    string           `json:"object"`
	Available map[string]int64 `json:"available"`
	Customer  string           `json:"customer"`
	Livemode  bool             `json:"livemode"`
}

type CustomerDiscount struct {
	Coupon   *Coupon `json:"coupon"`
	Customer string  `json:"customer"`
	End      int64   `json:"end"`
	Id       string  `json:"id"`
	Object   string  `json:"object"`
	Start    int64   `json:"start"`
}

type CustomerInvoiceSettings struct {
	CustomFields         []*InvoiceCustomField    `json:"custom_fields"`
	DefaultPaymentMethod *string                  `json:"default_payment_method"`
	Footer               string                   `json:"footer"`
	RenderingOptions     *InvoiceRenderingOptions `json:"rendering_options"`
}

type InvoiceCustomField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type InvoiceRenderingOptions struct {
	AmountTaxDisplay string `json:"amount_tax_display"`
}

type CustomerShipping struct {
	Address *CustomerAddress `json:"address"`
	Name    string           `json:"name"`
	Phone   string           `json:"phone"`
}

type CustomerSourceList struct {
	Object  string        `json:"object"`
	Data    []interface{} `json:"data"`
	HasMore bool          `json:"has_more"`
	Url     string        `json:"url"`
}

type SubscriptionList struct {
	Object  string        `json:"object"`
	Data    []interface{} `json:"data"`
	HasMore bool          `json:"has_more"`
	Url     string        `json:"url"`
}

type CustomerTax struct {
	AutomaticTax string               `json:"automatic_tax"`
	IpAddress    string               `json:"ip_address"`
	Location     *CustomerTaxLocation `json:"location"`
}

type CustomerTaxLocation struct {
	Country string `json:"country"`
	Source  string `json:"source"`
	State   string `json:"state"`
}

type CustomerTaxIdList struct {
	Object  string          `json:"object"`
	Data    []CustomerTaxId `json:"data"`
	HasMore bool            `json:"has_more"`
	Url     string          `json:"url"`
}

type CustomerTaxId struct {
	Id       string `json:"id"`
	Object   string `json:"object"`
	Country  string `json:"country"`
	Customer string `json:"customer"`
	Type     string `json:"type"`
	Value    string `json:"value"`
}

type TestClock struct {
	Id         string `json:"id"`
	Object     string `json:"object"`
	Created    int64  `json:"created"`
	DelayedBy  int64  `json:"delayed_by"`
	FrozenTime int64  `json:"frozen_time"`
	Livemode   bool   `json:"livemode"`
	Name       string `json:"name"`
	Status     string `json:"status"`
}

type Coupon struct {
	Id               string  `json:"id"`
	Object           string  `json:"object"`
	AmountOff        int64   `json:"amount_off"`
	Created          int64   `json:"created"`
	Currency         string  `json:"currency"`
	Duration         string  `json:"duration"`
	DurationInMonths int64   `json:"duration_in_months"`
	Livemode         bool    `json:"livemode"`
	MaxRedemptions   int64   `json:"max_redemptions"`
	Name             string  `json:"name"`
	PercentOff       float64 `json:"percent_off"`
	RedeemBy         int64   `json:"redeem_by"`
	TimesRedeemed    int64   `json:"times_redeemed"`
	Valid            bool    `json:"valid"`
}

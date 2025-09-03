package payment

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

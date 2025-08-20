package payment

// Setup Products
type SetupProductsReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
}

type SetupProductsResponse struct {
	OneTimeProduct struct {
		ID      string `json:"id"`
		PriceID string `json:"priceId"`
	} `json:"oneTimeProduct"`
	SubscriptionProduct struct {
		ID      string `json:"id"`
		PriceID string `json:"priceId"`
	} `json:"subscriptionProduct"`
}

// Create Customer

type CreateCustomerReq struct {
	Email string `json:"email"`
}

type CreateCustomerResponse struct {
	CustomerID string `json:"customerId"`
}

type CreatePaymentIntentRequest struct {
	Amount     int64  `json:"amount"` // Amount in cents
	CustomerID string `json:"customerId,omitempty"`
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

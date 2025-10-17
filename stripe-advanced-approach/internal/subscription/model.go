package subscription

// Request types
type SubscribeRequest struct {
	PriceID string `json:"price_id" binding:"required"`
}

// Response types
type SubscribeResponse struct {
	ClientSecret string `json:"client_secret,omitempty"`
	CheckoutURL  string `json:"checkout_url,omitempty"`
}

// Minimal response for frontend UI decisions
type SubscriptionStatusResponse struct {
	HasAccess         bool   `json:"has_access"`           // Simple boolean for "can see premium page?"
	Status            string `json:"status"`               // "active", "canceled", "past_due", "none"
	CancelAtPeriodEnd bool   `json:"cancel_at_period_end"` // Show "Renews on" vs "Expires on"
}

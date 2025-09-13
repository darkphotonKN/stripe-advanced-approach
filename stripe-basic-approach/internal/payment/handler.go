package payment

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

type Handler struct {
	service Service
}

type Service interface {
	SetupProducts(context.Context, *SetupProductsReq) (*SetupProductsResp, error)
	GetProducts(ctx context.Context) (*ProductListResponse, error)
	CreateCustomer(ctx context.Context, userId uuid.UUID, email string) (string, error)
	SaveCard(ctx context.Context, customerId string) (string, error)
	CreatePaymentIntent(ctx context.Context, amount int64, customerId string) (*CreatePaymentIntentResponse, error)
	PurchaseProduct(ctx context.Context, userId uuid.UUID, req *PurchaseProductRequest) (*PurchaseProductResponse, error)
	SetupSubscription(ctx context.Context, request *SetupProductsReq) (*SetupProductsResp, error)
	SubscribeToProduct(ctx context.Context, req *SubscribeRequest) (*SubscribeResponse, error)

	// flow based methods
	ProcessWebhookEvent(ctx context.Context, event *stripe.Event) error
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SetupProducts(c *gin.Context) {
	var request SetupProductsReq
	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Printf("\nError when parsing json: %+v\n\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.SetupProducts(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) CreateCustomer(c *gin.Context) {
	userIdStr, _ := c.Get("user_id")
	email, _ := c.Get("email")

	fmt.Printf("\nCreating customer with userId: %s and email: %s\n\n", userIdStr, email)

	userId, _ := uuid.Parse(userIdStr.(string))

	customerId, err := h.service.CreateCustomer(c.Request.Context(), userId, email.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"customer_id": customerId})
}

func (h *Handler) SaveCard(c *gin.Context) {
	var req struct {
		CustomerID string `json:"customer_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientSecret, err := h.service.SaveCard(c.Request.Context(), req.CustomerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"client_secret": clientSecret})
}

func (h *Handler) CreatePaymentIntent(c *gin.Context) {
	var req CreatePaymentIntentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// validation
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than 0"})
		return
	}
	if req.CustomerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer ID required"})
		return
	}

	result, err := h.service.CreatePaymentIntent(c.Request.Context(), req.Amount, req.CustomerID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment intent"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) SetupSubscription(c *gin.Context) {
	var request SetupProductsReq

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.SetupSubscription(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetProducts(c *gin.Context) {
	resp, err := h.service.GetProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) PurchaseProduct(c *gin.Context) {

	userIdStr, _ := c.Get("user_id")

	fmt.Printf("\nPurchasing product with userId: %s\n\n", userIdStr)

	userId, _ := uuid.Parse(userIdStr.(string))

	var req PurchaseProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.PurchaseProduct(c.Request.Context(), userId, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) SubscribeToProduct(c *gin.Context) {
	var req SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.SubscribeToProduct(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) HandleStripeWebhook(c *gin.Context) {
	// Read raw bytes instead of using gin's ShouldBindJSON because:
	// 1. Stripe's webhook signature is calculated from the exact bytes sent
	// 2. ShouldBindJSON would parse/reformat the JSON, breaking signature verification
	// 3. This ensures the webhook is authentic and hasn't been tampered with

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Printf("\nError reading body: %+v\n\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading body"})
		return
	}

	// Get Stripe signature from headers
	signature := c.GetHeader("Stripe-Signature")
	if signature == "" {
		fmt.Printf("\nError Missing Signature: %+v\n\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing signature"})
		return
	}

	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

	fmt.Printf("\nStripe webhook secret from ENV: %s\n\n", webhookSecret)
	event, err := webhook.ConstructEvent(body, signature, webhookSecret)
	if err != nil {
		fmt.Printf("\nError invalid signature: %+v\n\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		return
	}

	// validate signature in service method
	err = h.service.ProcessWebhookEvent(c.Request.Context(), &event)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to process stripe event.", "error": err})
		return
	}
}

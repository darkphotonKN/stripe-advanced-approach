package payment

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

type Service interface {
	SetupProducts(context.Context, *SetupProductsReq) (*SetupProductsResp, error)
	CreateCustomer(ctx context.Context, userId uuid.UUID, email string) (string, error)
	SaveCard(ctx context.Context, customerId string) (string, error)
	CreatePaymentIntent(ctx context.Context, amount int64, customerId string) (string, error)
	CreateSubscription(ctx context.Context, priceId, customerId, email string) (*SubscriptionResp, error)
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SetupProducts(c *gin.Context) {
	fmt.Println("Creating Stripe Product...")

	var request SetupProductsReq
	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Printf("\nError when parsing json: %+v\n\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("create product request:", request)

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

	clientSecret, err := h.service.CreatePaymentIntent(c.Request.Context(), req.Amount, req.CustomerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment intent"})
		return
	}

	c.JSON(http.StatusOK, CreatePaymentIntentResponse{
		ClientSecret: clientSecret,
	})
}

func (h *Handler) CreateSubscription(c *gin.Context) {
	var req struct {
		PriceID    string `json:"price_id" binding:"required"`
		CustomerID string `json:"customer_id" binding:"required"`
		Email      string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.CreateSubscription(c.Request.Context(), req.PriceID, req.CustomerID, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

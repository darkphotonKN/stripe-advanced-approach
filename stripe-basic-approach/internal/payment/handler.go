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
	SetupProducts(context.Context, *SetupProductsReq) error
	CreateCustomer(ctx context.Context, userId uuid.UUID, email string) error
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

	if err := h.service.SetupProducts(c.Request.Context(), &request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, request)
}

func (h *Handler) CreateCustomer(c *gin.Context) {
	userIdStr, _ := c.Get("user_id")
	email, _ := c.Get("email")

	fmt.Printf("\nCreating customer with userId: %s and email: %s\n\n", userIdStr, email)

	userId, _ := uuid.Parse(userIdStr.(string))

	if err := h.service.CreateCustomer(c.Request.Context(), userId, email.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Customer successfully created."})
}

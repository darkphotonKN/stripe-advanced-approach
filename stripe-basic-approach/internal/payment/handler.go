package payment

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

type Service interface {
	SetupProducts(context.Context, *SetupProductsReq) error
	CreateCustomer(context.Context, *CreateCustomerReq) error
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SetupProducts(c *gin.Context) {
	var request SetupProductsReq
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SetupProducts(c.Request.Context(), &request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, request)
}

func (h *Handler) CreateCustomer(c *gin.Context) {
	var request CreateCustomerReq
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateCustomer(c.Request.Context(), &request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, request)
}

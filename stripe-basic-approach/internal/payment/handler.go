package payment

import (
	"github.com/docker/distribution/uuid"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

type Service interface {
	CreateCustomer(userId uuid.UUID) error
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateCustomer(c *gin.Context) {
	// var product Product
	// if err := c.ShouldBindJSON(&product); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }
	//
	// if err := h.service.CreateCustomer(c.Request.Context(), &product); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	//
	// c.JSON(http.StatusCreated, product)
}

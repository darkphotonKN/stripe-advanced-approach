package config

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/darkphotonKN/stripe-basic-approach/internal/middleware"
	"github.com/darkphotonKN/stripe-basic-approach/internal/payment"
	"github.com/darkphotonKN/stripe-basic-approach/internal/product"
	"github.com/darkphotonKN/stripe-basic-approach/internal/user"
)

func SetupRoutes(db *sqlx.DB) *gin.Engine {
	router := gin.Default()

	userRepo := user.NewRepository(db)
	productRepo := product.NewRepository(db)

	userService := user.NewService(userRepo)
	productService := product.NewService(productRepo)

	userHandler := user.NewHandler(userService)
	productHandler := product.NewHandler(productService)

	// stripeProcessor := payment.NewStripeProcessor()

	api := router.Group("/api")

	api.POST("/signup", userHandler.SignUp)
	api.POST("/signin", userHandler.SignIn)

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())

	protected.GET("/users", userHandler.List)
	protected.GET("/users/:id", userHandler.Get)
	protected.PUT("/users/:id", userHandler.Update)
	protected.DELETE("/users/:id", userHandler.Delete)

	protected.GET("/products", productHandler.List)
	protected.POST("/products", productHandler.Create)
	protected.GET("/products/:id", productHandler.Get)
	protected.PUT("/products/:id", productHandler.Update)
	protected.DELETE("/products/:id", productHandler.Delete)

	// payment setup
	paymentService := payment.NewService()
	paymentHandler := payment.NewHandler(paymentService)

	protected.POST("/create-customer", paymentHandler.CreateCustomer)

	// router.HandleFunc("/save-card", h.SaveCard).Methods("POST")

	return router
}


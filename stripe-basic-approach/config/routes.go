package config

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/darkphotonKN/stripe-basic-approach/internal/middleware"
	"github.com/darkphotonKN/stripe-basic-approach/internal/payment"
	"github.com/darkphotonKN/stripe-basic-approach/internal/user"
)

func SetupRoutes(db *sqlx.DB) *gin.Engine {
	router := gin.Default()

	// NOTE: debugging middleware
	router.Use(func(c *gin.Context) {
		fmt.Println("Incoming request to:", c.Request.Method, c.Request.URL.Path, "from", c.Request.Host)
		c.Next()
	})

	// TODO: CORS for development, remove in PROD
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

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

	// payment setup
	paymentService := payment.NewService()
	paymentHandler := payment.NewHandler(paymentService)

	protected.POST("/setup-products", paymentHandler.SetupProducts)
	protected.POST("/create-customer", paymentHandler.CreateCustomer)
	protected.POST("/save-card", paymentHandler.SaveCard)
	protected.POST("/create-payment-intent", paymentHandler.CreatePaymentIntent)
	protected.POST("/create-subscription", paymentHandler.CreateSubscription)

	return router
}

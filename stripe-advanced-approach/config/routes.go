package config

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/darkphotonKN/stripe-advanced-approach/internal/interfaces"
	"github.com/darkphotonKN/stripe-advanced-approach/internal/middleware"
	"github.com/darkphotonKN/stripe-advanced-approach/internal/payment"
	"github.com/darkphotonKN/stripe-advanced-approach/internal/user"
)

func SetupRoutes(db *sqlx.DB, cacheClient interfaces.Cache) *gin.Engine {
	router := gin.Default()

	// NOTE: debugging middleware
	router.Use(func(c *gin.Context) {
		fmt.Println("Incoming request to:", c.Request.Method, c.Request.URL.Path, "from", c.Request.Host)
		c.Next()
	})

	// TODO: CORS for development, remove in PROD
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Content-Type", "Authorization",
			"Accept",
			"X-Requested-With",
			"X-CSRF-Token",
			"Cache-Control",
			"*", // TODO: remove in prod
		},
		AllowCredentials: true,
	}))

	api := router.Group("/api")

	// user setup
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)
	api.POST("/signup", userHandler.SignUp)
	api.POST("/signin", userHandler.SignIn)

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())

	protected.GET("/users", userHandler.List)
	protected.GET("/users/stripe-customer", userHandler.GetStripeCustomer)
	protected.GET("/users/:id", userHandler.Get)
	protected.PUT("/users/:id", userHandler.Update)
	protected.DELETE("/users/:id", userHandler.Delete)

	// payment setup
	stripeProcessor := payment.NewStripeProcessor()
	paymentRepository := payment.NewRepository(db)
	paymentService := payment.NewService(paymentRepository, userService, stripeProcessor, cacheClient)

	// injecting proper payment service after completing payment service initialization
	userService.SetPaymentService(paymentService)

	paymentHandler := payment.NewHandler(paymentService)

	// for stripe webhooks
	stripeWebhookAPI := router.Group("/")
	stripeWebhookAPI.POST("/webhook/stripe", paymentHandler.HandleStripeWebhook)

	// payment service endpoints
	paymentRoutes := protected.Group("/payment")
	paymentRoutes.POST("/setup-products", paymentHandler.SetupProducts)
	paymentRoutes.POST("/setup-subscription", paymentHandler.SetupSubscription)
	paymentRoutes.GET("/products", paymentHandler.GetProducts)
	paymentRoutes.POST("/create-customer", paymentHandler.CreateCustomer)
	paymentRoutes.POST("/save-card", paymentHandler.SaveCard)
	paymentRoutes.POST("/create-payment-intent", paymentHandler.CreatePaymentIntent)
	paymentRoutes.POST("/purchase-product", paymentHandler.PurchaseProduct)
	paymentRoutes.POST("/subscribe-to-product", paymentHandler.SubscribeToProduct)

	// subscription endpoints (part of payment service)
	paymentRoutes.POST("/subscription/subscribe", paymentHandler.Subscribe)
	paymentRoutes.GET("/subscription/status", paymentHandler.GetSubscriptionStatus)

	return router
}

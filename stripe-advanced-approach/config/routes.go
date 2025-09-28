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

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	api := router.Group("/api")

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
	paymentService := payment.NewService(paymentRepository, userService, stripeProcessor)
	paymentHandler := payment.NewHandler(paymentService)

	// for stripe webhooks
	stripeWebhookAPI := router.Group("/")
	stripeWebhookAPI.POST("/webhook/stripe", paymentHandler.HandleStripeWebhook)

	protected.POST("/setup-products", paymentHandler.SetupProducts)
	protected.POST("/setup-subscription", paymentHandler.SetupSubscription)
	protected.GET("/products", paymentHandler.GetProducts)
	protected.POST("/create-customer", paymentHandler.CreateCustomer)
	protected.POST("/save-card", paymentHandler.SaveCard)
	protected.POST("/create-payment-intent", paymentHandler.CreatePaymentIntent)
	protected.POST("/purchase-product", paymentHandler.PurchaseProduct)
	protected.POST("/subscribe-to-product", paymentHandler.SubscribeToProduct)

	return router
}

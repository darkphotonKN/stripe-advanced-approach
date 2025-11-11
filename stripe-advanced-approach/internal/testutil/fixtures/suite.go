package fixtures

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/darkphotonKN/stripe-advanced-approach/internal/payment"
	"github.com/darkphotonKN/stripe-advanced-approach/internal/redis"
	"github.com/darkphotonKN/stripe-advanced-approach/internal/user"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stripe/stripe-go/v82"
)

// TestSuite holds all application services with proper dependency injection
type TestSuite struct {
	Ctx         context.Context
	DB          *sqlx.DB
	RedisClient *redis.Client
	CleanupFunc func()

	// Services (fully configured with DI)
	UserService    *user.Service
	PaymentService *payment.Service

	// Repositories
	UserRepo    user.Repository
	PaymentRepo payment.Repository

	// Handlers
	UserHandler    *user.Handler
	PaymentHandler *payment.Handler

	// Test Data
	TestUser TestUser
}

type TestUser struct {
	UserID     uuid.UUID
	CustomerID string
	Email      string
	Password   string
}

// SetupTestSuite creates fully configured test environment with all services
func SetupTestSuite(t *testing.T) *TestSuite {
	t.Helper()

	// Load environment variables
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Log("No .env file found, using environment variables")
	}

	// Setup Stripe
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	if stripe.Key == "" {
		t.Skip("STRIPE_SECRET_KEY not set, skipping test")
	}

	// Setup database connection
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	// Setup Redis client
	redisClient := redis.NewClient()

	// Setup repositories
	userRepo := user.NewRepository(db)
	paymentRepo := payment.NewRepository(db)

	// Setup services with proper dependency injection
	// Step 1: Create user service first (without payment service)
	userService := user.NewService(userRepo)

	// Step 2: Create payment service (can accept user service)
	stripeProcessor := payment.NewStripeProcessor()
	paymentService := payment.NewService(paymentRepo, userService, stripeProcessor, redisClient)

	// Step 3: Inject payment service back into user service (resolves circular dependency)
	userService.SetPaymentService(paymentService)

	// Setup handlers
	userHandler := user.NewHandler(userService)
	paymentHandler := payment.NewHandler(paymentService)

	// Create test user data
	testUser := TestUser{
		UserID:     uuid.New(),
		CustomerID: "cus_test_" + uuid.New().String()[:8],
		Email:      "nov7subscriber@test.com",
		Password:   "123456",
	}

	// Cleanup function
	cleanupFunc := func() {
		db.Close()
	}

	return &TestSuite{
		Ctx:         context.Background(),
		DB:          db,
		RedisClient: redisClient,
		CleanupFunc: cleanupFunc,

		UserService:    userService,
		PaymentService: paymentService,

		UserRepo:    userRepo,
		PaymentRepo: paymentRepo,

		UserHandler:    userHandler,
		PaymentHandler: paymentHandler,

		TestUser: testUser,
	}
}

// SetupTestSuiteWithCustomUser creates test suite with custom user data
func SetupTestSuiteWithCustomUser(t *testing.T, customUser TestUser) *TestSuite {
	suite := SetupTestSuite(t)
	suite.TestUser = customUser
	return suite
}
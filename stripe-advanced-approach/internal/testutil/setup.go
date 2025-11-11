package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/darkphotonKN/stripe-advanced-approach/internal/redis"
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

	// Services and repositories - all as interfaces to avoid import cycles
	UserService    interface{}
	PaymentService interface{}

	// Repositories
	UserRepo    interface{}
	PaymentRepo interface{}

	// Handlers
	UserHandler    interface{}
	PaymentHandler interface{}

	// Test Data
	TestUser TestUser
}

type TestUser struct {
	UserID     uuid.UUID
	CustomerID string
	Email      string
	Password   string
}

// SetupBasicTestSuite creates basic test infrastructure (DB, Redis)
func SetupBasicTestSuite(t *testing.T) *TestSuite {
	t.Helper()

	// Load environment variables
	if err := godotenv.Load("../../.env"); err != nil {
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

	// Payment repository will be created in individual tests

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

		TestUser: testUser,
	}
}

// SetupBasicTestSuiteWithCustomUser creates test suite with custom user data
func SetupBasicTestSuiteWithCustomUser(t *testing.T, customUser TestUser) *TestSuite {
	suite := SetupBasicTestSuite(t)
	suite.TestUser = customUser
	return suite
}


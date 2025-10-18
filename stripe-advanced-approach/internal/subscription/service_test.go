package subscription

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

type SubscriptionTestSuite struct {
	ctx            context.Context
	service        Service
	redisClient    *redis.Client
	userService    user.Service
	paymentService payment.Service
	db             *sqlx.DB
	cleanupFunc    func()
}

func setupTestSuite(t *testing.T) *SubscriptionTestSuite {
	t.Helper()

	if err := godotenv.Load("../../.env"); err != nil {
		t.Log("No .env file found, using environment variables")
	}

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

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		t.Fatal("Failed to connect to database:", err)
	}

	redisClient := redis.NewClient()

	ctx := context.Background()
	if err := redisClient.Connect(ctx); err != nil {
		t.Fatal("Failed to connect to Redis:", err)
	}

	// Setup dependencies
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)

	paymentRepo := payment.NewRepository(db)
	stripeProcessor := payment.NewStripeProcessor()
	paymentService := payment.NewService(paymentRepo, userService, stripeProcessor, redisClient)

	// Inject payment service into user service
	userService.SetPaymentService(paymentService)

	// Setup subscription service
	subscriptionService := NewService(userService, paymentService)

	cleanup := func() {
		if err := db.Close(); err != nil {
			t.Logf("Failed to close database connection: %v", err)
		}
		if err := redisClient.Close(); err != nil {
			t.Logf("Failed to close Redis connection: %v", err)
		}
	}

	return &SubscriptionTestSuite{
		ctx:            ctx,
		service:        subscriptionService,
		redisClient:    redisClient,
		userService:    userService,
		paymentService: paymentService,
		db:             db,
		cleanupFunc:    cleanup,
	}
}

func (suite *SubscriptionTestSuite) Cleanup() {
	if suite.cleanupFunc != nil {
		suite.cleanupFunc()
	}
}

func TestGetSubscriptionStatus(t *testing.T) {
	suite := setupTestSuite(t)
	defer suite.Cleanup()

	// userId for testing
	userId := uuid.MustParse("158c0a1e-583e-4bb6-9378-d980a17c5809")

	// Call GetSubscriptionStatus
	status, err := suite.service.GetSubscriptionStatus(suite.ctx, userId)
	if err != nil {
		t.Logf("Error when getting subscription status: %v", err)
	}

	if status != nil {
		t.Logf("Subscription status: has_access=%v, status=%s, cancel_at_period_end=%v",
			status.HasAccess, status.Status, status.CancelAtPeriodEnd)
	} else {
		t.Log("Subscription status returned nil")
	}
}

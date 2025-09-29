package payment

import (
	"context"
	"os"
	"testing"

	"github.com/darkphotonKN/stripe-advanced-approach/internal/redis"
	"github.com/darkphotonKN/stripe-advanced-approach/internal/user"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v82"
)

type PaymentTestSuite struct {
	ctx             context.Context
	service         *service
	redisClient     *redis.Client
	repo            Repository
	userService     PaymentUserService
	stripeProcessor PaymentProcessor
	cleanupFunc     func()
}

func setupTestSuite(t *testing.T) *PaymentTestSuite {
	t.Helper()

	if err := godotenv.Load("../../.env"); err != nil {
		t.Log("No .env file found, using environment variables")
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	if stripe.Key == "" {
		t.Skip("STRIPE_SECRET_KEY not set, skipping test")
	}

	redisClient := redis.NewClient()

	ctx := context.Background()
	if err := redisClient.Connect(ctx); err != nil {
		t.Fatal("Failed to connect to Redis:", err)
	}

	emptyTestRepo := NewRepository(&sqlx.DB{})
	emptyTestUserRepo := user.NewRepository(&sqlx.DB{})
	emptyTestUserServ := user.NewService(emptyTestUserRepo)
	emptyTestStripeProcessor := NewStripeProcessor()

	service := NewService(emptyTestRepo, emptyTestUserServ, emptyTestStripeProcessor, redisClient)

	cleanup := func() {
		if err := redisClient.Close(); err != nil {
			t.Logf("Failed to close Redis connection: %v", err)
		}
	}

	return &PaymentTestSuite{
		ctx:             ctx,
		service:         service,
		redisClient:     redisClient,
		repo:            emptyTestRepo,
		userService:     emptyTestUserServ,
		stripeProcessor: emptyTestStripeProcessor,
		cleanupFunc:     cleanup,
	}
}

func (suite *PaymentTestSuite) Cleanup() {
	if suite.cleanupFunc != nil {
		suite.cleanupFunc()
	}
}

func TestSyncStripeDataToStorage(t *testing.T) {
	suite := setupTestSuite(t)
	defer suite.Cleanup()

	err := suite.service.SyncStripeDataToStorage(suite.ctx, "cus_Sv7wygVCW12KLE")
	if err != nil {
		t.Logf("errored when attempting to sync stripe data to storage: %v", err)
	}
}

func TestGetStripeData(t *testing.T) {
	suite := setupTestSuite(t)
	defer suite.Cleanup()

	err := suite.service.SyncStripeDataToStorage(suite.ctx, "cus_Sv7wygVCW12KLE")
	if err != nil {
		t.Logf("Failed to sync data: %v", err)
	}

	err = suite.service.GetStripeData(suite.ctx, "stripe:customer:cus_Sv7wygVCW12KLE")
	if err != nil {
		t.Logf("Failed to get cached data: %v", err)
	}
}


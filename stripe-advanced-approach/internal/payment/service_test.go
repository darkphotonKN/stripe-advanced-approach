package payment

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/darkphotonKN/stripe-advanced-approach/internal/redis"
	"github.com/darkphotonKN/stripe-advanced-approach/internal/user"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stripe/stripe-go/v82"
)

type PaymentTestSuite struct {
	ctx             context.Context
	service         *service
	redisClient     *redis.Client
	repo            Repository
	userService     PaymentUserService
	stripeProcessor PaymentProcessor
	db              *sqlx.DB
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

	testRepo := NewRepository(db)
	testUserRepo := user.NewRepository(db)
	testUserServ := user.NewService(testUserRepo)
	testStripeProcessor := NewStripeProcessor()

	service := NewService(testRepo, testUserServ, testStripeProcessor, redisClient)

	cleanup := func() {
		if err := db.Close(); err != nil {
			t.Logf("Failed to close database connection: %v", err)
		}
		if err := redisClient.Close(); err != nil {
			t.Logf("Failed to close Redis connection: %v", err)
		}
	}

	return &PaymentTestSuite{
		ctx:             ctx,
		service:         service,
		redisClient:     redisClient,
		repo:            testRepo,
		userService:     testUserServ,
		stripeProcessor: testStripeProcessor,
		db:              db,
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

	err := suite.service.SyncStripeDataToStorage(suite.ctx, "cus_TJ60h8lizGx9CV")
	if err != nil {
		t.Logf("errored when attempting to sync stripe data to storage: %v", err)
	}
}

func TestGetStripeData(t *testing.T) {
	suite := setupTestSuite(t)
	defer suite.Cleanup()

	_, err := suite.service.GetStripeData(suite.ctx, "cus_TJ60h8lizGx9CV")
	if err != nil {
		t.Logf("Failed to get cached data: %v", err)
	}
}

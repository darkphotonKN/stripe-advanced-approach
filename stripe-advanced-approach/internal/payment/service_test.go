package payment

import (
	"context"
	"testing"

	"github.com/darkphotonKN/stripe-advanced-approach/internal/testutil"
	"github.com/google/uuid"
)

// Mock user service for testing
type mockUserService struct{}

func (m *mockUserService) GetByID(ctx context.Context, id uuid.UUID) (interface{}, error) {
	return nil, nil
}

func (m *mockUserService) GetByStripeCustomerID(ctx context.Context, stripeCustomerID string) (interface{}, error) {
	return nil, nil
}

func (m *mockUserService) UpdateStripeCustomer(ctx context.Context, userId uuid.UUID, customerId string) error {
	return nil
}

type PaymentTestSuite struct {
	ctx             context.Context
	service         *service
	redisClient     *redis.Client
	repo            Repository
	userService     PaymentUserService
	stripeProcessor PaymentProcessor
	db              *sqlx.DB
	cleanupFunc     func()

	// metadata
	testUser TestUser
}

type TestUser struct {
	userId     uuid.UUID
	customerId string
}

// setupTestSuite creates a legacy PaymentTestSuite for backward compatibility
func setupTestSuite(t *testing.T, testUser TestUser) *PaymentTestSuite {
	// Create custom test user for shared setup
	customTestUser := testutil.TestUser{
		UserID:     testUser.userId,
		CustomerID: testUser.customerId,
		Email:      "test@example.com",
		Password:   "password123",
	}

	sharedSuite := testutil.SetupBasicTestSuiteWithCustomUser(t, customTestUser)

	// Setup payment-specific services
	repo := NewRepository(sharedSuite.DB)
	stripeProcessor := NewStripeProcessor()

	// Create mock user service
	mockUserService := &mockUserService{}

	// Create payment service
	paymentService := NewService(repo, mockUserService, stripeProcessor, sharedSuite.RedisClient)

	// Create legacy suite for backward compatibility
	return &PaymentTestSuite{
		ctx:             sharedSuite.Ctx,
		service:         paymentService,
		redisClient:     sharedSuite.RedisClient,
		repo:            repo,
		userService:     mockUserService,
		stripeProcessor: stripeProcessor,
		db:              sharedSuite.DB,
		cleanupFunc:     sharedSuite.CleanupFunc,
		testUser:        testUser,
	}
}

func (suite *PaymentTestSuite) Cleanup() {
	if suite.cleanupFunc != nil {
		suite.cleanupFunc()
	}
}

// Test user meta data
var testUser = TestUser{
	userId:     uuid.MustParse("7fdf18f8-78b6-4ca8-b2dd-d6dfb8286fe7"),
	customerId: "cus_TNYXOoNL2FlSGA",
}

// Tests

func TestSyncStripeDataToStorage(t *testing.T) {
	suite := setupTestSuite(t, testUser)
	defer suite.Cleanup()

	err := suite.service.SyncStripeDataToStorage(suite.ctx, suite.testUser.customerId)
	if err != nil {
		t.Logf("errored when attempting to sync stripe data to storage: %v", err)
	}
}

func TestGetStripeData(t *testing.T) {
	suite := setupTestSuite(t, testUser)
	defer suite.Cleanup()

	_, err := suite.service.GetStripeData(suite.ctx, suite.testUser.customerId)
	if err != nil {
		t.Logf("Failed to get cached data: %v", err)
	}
}

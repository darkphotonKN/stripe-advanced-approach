package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darkphotonKN/stripe-advanced-approach/internal/testutil"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock payment service for testing
type mockPaymentService struct{}

func (m *mockPaymentService) AddCacheUserIdToCusId(ctx context.Context, userId uuid.UUID, customerId string) error {
	return nil
}

func (m *mockPaymentService) AddCacheCusIdToUserId(ctx context.Context, customerId string, userId uuid.UUID) error {
	return nil
}

func (m *mockPaymentService) CreateCustomer(ctx context.Context, userId uuid.UUID, email string) (string, error) {
	return "cus_test_mock", nil
}

func (m *mockPaymentService) SyncStripeDataToStorage(ctx context.Context, customerId string) error {
	return nil
}

// TestSignIn tests the signin functionality with full service integration
func TestSignIn(t *testing.T) {
	suite := testutil.SetupBasicTestSuite(t)
	defer suite.CleanupFunc()

	// Setup services with DI manually to avoid import cycles
	userRepo := NewRepository(suite.DB)
	userService := NewService(userRepo)

	// For testing, create mock payment service interface
	mockPaymentService := &mockPaymentService{}
	userService.SetPaymentService(mockPaymentService)

	userHandler := NewHandler(userService)

	// Setup test router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/signin", userHandler.SignIn)

	// Use test user from suite
	testEmail := suite.TestUser.Email
	testPassword := suite.TestUser.Password

	// Prepare signin request
	requestBody := SignInRequest{
		Email:    testEmail,
		Password: testPassword,
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Record response
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Check status code
	assert.Equal(t, http.StatusOK, w.Code, "Expected successful signin")

	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Validate response
	assert.NotEmpty(t, response["access_token"], "Access token should be present")
	assert.NotNil(t, response["user"], "User object should be returned")

	// Check user email in response
	userData := response["user"].(map[string]interface{})
	assert.Equal(t, testEmail, userData["email"], "Email should match")

	t.Logf("SignIn test passed for user: %s", testEmail)
}

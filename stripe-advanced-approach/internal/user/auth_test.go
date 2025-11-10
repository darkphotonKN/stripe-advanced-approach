package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/darkphotonKN/stripe-advanced-approach/internal/redis"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test SignIn functionality
func TestSignIn(t *testing.T) {
	// Load environment variables
	if err := godotenv.Load("../../.env"); err != nil {
		t.Log("No .env file found, using environment variables")
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
	defer db.Close()

	// Setup repository
	repo := NewRepository(db)

	// Setup service
	service := NewService(repo)

	// Setup handler
	handler := NewHandler(service)

	// Setup test router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/signin", handler.SignIn)

	// Test user credentials
	testEmail := "nov7subscriber@test.com"
	testPassword := "123456"

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

// Mock payment service for testing
type mockPaymentService struct{}

func (m *mockPaymentService) CreateCustomer(ctx context.Context, userId uuid.UUID, email string) (string, error) {
	return "cus_mock123", nil
}


package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/darkphotonKN/stripe-advanced-approach/internal/testutil/fixtures"
	"github.com/darkphotonKN/stripe-advanced-approach/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSignIn tests the signin functionality with full service integration
func TestSignIn(t *testing.T) {
	suite := fixtures.SetupTestSuite(t)
	defer suite.CleanupFunc()

	// Setup test router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/signin", suite.UserHandler.SignIn)

	// Use test user from suite
	testEmail := suite.TestUser.Email
	testPassword := suite.TestUser.Password

	// Prepare signin request
	requestBody := user.SignInRequest{
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

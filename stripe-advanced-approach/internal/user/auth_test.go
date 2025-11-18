package user_test

import (
	"testing"

	"github.com/darkphotonKN/stripe-advanced-approach/internal/testutil"
	"github.com/darkphotonKN/stripe-advanced-approach/internal/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSignIn tests the signin functionality with service integration
func TestSignIn(t *testing.T) {
	// Setup full test suite with all services configured
	suite := testutil.SetupFull(t)
	defer suite.CleanupFunc()

	// Create test user first
	testUser := &user.User{
		Email:    suite.TestUser.Email,
		Password: suite.TestUser.Password,
		Name:     "Test User",
	}

	err := suite.UserService.Create(suite.Ctx, testUser)
	require.NoError(t, err, "Failed to create test user")

	// Test authentication
	authenticatedUser, err := suite.UserService.Authenticate(suite.Ctx, suite.TestUser.Email, suite.TestUser.Password)

	// Assertions
	assert.NoError(t, err, "Authentication should succeed")
	assert.NotNil(t, authenticatedUser, "User should be returned")
	assert.Equal(t, suite.TestUser.Email, authenticatedUser.Email, "Email should match")
	assert.Empty(t, authenticatedUser.Password, "Password should be cleared")
	assert.NotNil(t, authenticatedUser.StripeCustomerID, "Stripe customer should be created")

	t.Logf("SignIn test passed for user: %s with Stripe customer: %s",
		authenticatedUser.Email, *authenticatedUser.StripeCustomerID)
}

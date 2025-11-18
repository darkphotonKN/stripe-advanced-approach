package payment_test

import (
	"testing"

	"github.com/darkphotonKN/stripe-advanced-approach/internal/testutil"
	"github.com/google/uuid"
)

// Test user meta data
var testUserData = testutil.TestUser{
	UserID:     uuid.MustParse("7fdf18f8-78b6-4ca8-b2dd-d6dfb8286fe7"),
	CustomerID: "cus_TNYXOoNL2FlSGA",
	Email:      "test@example.com",
	Password:   "password123",
}

// TestSaveCard tests the SaveCard functionality which internally uses SyncStripeDataToStorage
func TestSaveCard(t *testing.T) {
	// Setup full test suite with custom user data
	suite := testutil.SetupFullWithUser(t, testUserData)
	defer suite.CleanupFunc()

	// Test SaveCard which internally calls SyncStripeDataToStorage
	clientSecret, err := suite.PaymentService.SaveCard(suite.Ctx, testUserData.CustomerID)

	if err != nil {
		t.Logf("Error saving card (expected for test customer): %v", err)
	} else {
		t.Logf("Got client secret: %s", clientSecret)
	}
}

// TestGetSubscriptionStatus tests getting subscription status which may use cached data
func TestGetSubscriptionStatus(t *testing.T) {
	// Setup full test suite with custom user data
	suite := testutil.SetupFullWithUser(t, testUserData)
	defer suite.CleanupFunc()

	// Test getting subscription status
	status, err := suite.PaymentService.GetSubscriptionStatus(suite.Ctx, testUserData.UserID)

	if err != nil {
		t.Logf("Error getting subscription status: %v", err)
	} else {
		t.Logf("Subscription status: %+v", status)
	}
}


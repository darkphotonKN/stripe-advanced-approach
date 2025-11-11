package payment_test

import (
	"testing"

	"github.com/darkphotonKN/stripe-advanced-approach/internal/testutil/fixtures"
	"github.com/google/uuid"
)

// Test user meta data
var testUser = fixtures.TestUser{
	UserID:     uuid.MustParse("7fdf18f8-78b6-4ca8-b2dd-d6dfb8286fe7"),
	CustomerID: "cus_TNYXOoNL2FlSGA",
	Email:      "test@example.com",
	Password:   "password123",
}

// Tests

func TestSyncStripeDataToStorage(t *testing.T) {
	suite := fixtures.SetupTestSuiteWithCustomUser(t, testUser)
	defer suite.CleanupFunc()

	err := suite.PaymentService.SyncStripeDataToStorage(suite.Ctx, testUser.CustomerID)
	if err != nil {
		t.Logf("errored when attempting to sync stripe data to storage: %v", err)
	}
}

func TestGetStripeData(t *testing.T) {
	suite := fixtures.SetupTestSuiteWithCustomUser(t, testUser)
	defer suite.CleanupFunc()

	_, err := suite.PaymentService.GetStripeData(suite.Ctx, testUser.CustomerID)
	if err != nil {
		t.Logf("Failed to get cached data: %v", err)
	}
}

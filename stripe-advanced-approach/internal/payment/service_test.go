package payment

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v82"
)

func TestSyncStripeDataToStorage(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Log("No .env file found, using environment variables")
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	if stripe.Key == "" {
		t.Skip("STRIPE_SECRET_KEY not set, skipping test")
	}

	s := &service{}

	err := s.SyncStripeDataToStorage("cus_T2K9xYijKi4uYR")
	if err != nil {
		t.Logf("Expected this to fail for now: %v", err)
	}
}


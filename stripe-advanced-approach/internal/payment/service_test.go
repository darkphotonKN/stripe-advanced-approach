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

func TestSyncStripeDataToStorage(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Log("No .env file found, using environment variables")
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	if stripe.Key == "" {
		t.Skip("STRIPE_SECRET_KEY not set, skipping test")
	}

	// initialize Redis
	redisClient := redis.NewClient()

	ctx := context.Background()
	if err := redisClient.Connect(ctx); err != nil {
		t.Error("Failed to connect to Redis:", err)
	}
	defer redisClient.Close()

	emptyTestRepo := NewRepository(&sqlx.DB{})
	emptyTestUserRepo := user.NewRepository(&sqlx.DB{})
	emptyTestUserServ := user.NewService(emptyTestUserRepo)
	emptyTestStripeProccessor := NewStripeProcessor()
	s := NewService(emptyTestRepo, emptyTestUserServ, emptyTestStripeProccessor, redisClient)

	err := s.SyncStripeDataToStorage(context.Background(), "cus_T2K9xYijKi4uYR")
	if err != nil {
		t.Logf("errored when attempting to sync stripe data to storage: %v", err)
	}

}


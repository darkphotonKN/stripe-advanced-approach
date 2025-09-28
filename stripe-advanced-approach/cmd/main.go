package main

import (
	"context"
	"log"
	"os"

	"github.com/darkphotonKN/stripe-advanced-approach/config"
	"github.com/darkphotonKN/stripe-advanced-approach/internal/redis"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v82"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := config.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// initialize Redis
	redisClient := redis.NewClient()

	ctx := context.Background()
	if err := redisClient.Connect(ctx); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer redisClient.Close()

	// test Redis connection
	if err := redisClient.Ping(ctx); err != nil {
		log.Fatal("Failed to ping Redis:", err)
	}

	// setup stripe
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// setup routes
	router := config.SetupRoutes(db, redisClient)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

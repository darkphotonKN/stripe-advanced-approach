package main

import (
	"log"
	"os"

	"github.com/darkphotonKN/stripe-advanced-approach/config"
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

	// setup stripe
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// setup routes
	router := config.SetupRoutes(db)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}


package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tryoasnafi/be-assignment/payment/internal/auth"
	"github.com/tryoasnafi/be-assignment/payment/internal/cors"
	"github.com/tryoasnafi/be-assignment/payment/internal/database"
	"github.com/tryoasnafi/be-assignment/payment/internal/transaction"
)

func init() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file", err)
	}
	// JWKs use to get the public key JWT Token
	go auth.FetchJWKs()
}

func main() {
	// Initialize database connection
	db, err := database.New()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors.Config())
	// router.Use(auth.Verify())

	// Misc endpoint
	router.GET(
		"/transaction-migrate",
		database.ValidateKey(),
		database.MigrationHandler,
	)

	// Register all routes
	apiRoute := router.Group("api")

	// Register all services
	// Transaction service
	transactionRepo := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepo)
	transaction.SetHandlers(apiRoute, transactionService)

	// starting the server
	addr := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start account service: %v", err)
	}
}

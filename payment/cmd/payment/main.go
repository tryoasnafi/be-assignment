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

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/tryoasnafi/be-assignment/payment/docs"
)

func init() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file", err)
	}
	// JWKs use to get the public key JWT Token
	go auth.FetchJWKs()
}

// @title           Payment Service API
// @version         1.0
// @description     This is a payment service - corebank.

// @host      localhost:9091
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token (Get JWT from header signin st-access-token).

// @externalDocs.description  User Auth API docs
// @externalDocs.url          https://localhost:9090/docs/index.html
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

	// Register all routes
	apiRoute := router.Group("api")

	// Register all services
	// Transaction service
	transactionRepo := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepo)
	transaction.SetHandlers(apiRoute, transactionService)

	// Misc endpoint
	apiRoute.POST("/transaction-migrate", database.ValidateKey(), database.MigrationHandler)
	apiRoute.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// starting the server
	addr := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start account service: %v", err)
	}
}

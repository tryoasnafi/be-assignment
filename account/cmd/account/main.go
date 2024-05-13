package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tryoasnafi/be-assignment/account/internal/account"
	"github.com/tryoasnafi/be-assignment/account/internal/auth"
	auth_supertokens "github.com/tryoasnafi/be-assignment/account/internal/auth-supertokens"
	"github.com/tryoasnafi/be-assignment/account/internal/cors"
	"github.com/tryoasnafi/be-assignment/account/internal/database"
	"github.com/tryoasnafi/be-assignment/account/internal/user"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/tryoasnafi/be-assignment/account/docs"
)

func init() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file", err)
	}
}

// @title           Account Service API
// @version         1.0
// @description     This is a account service - corebank.

// @host      localhost:9090
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token (Get JWT from header signin st-access-token).

// @externalDocs.description  Supertokens Auth
// @externalDocs.url          https://app.swaggerhub.com/apis/supertokens/FDI
func main() {
	// Initialize database connection
	db, err := database.New()
	if err != nil {
		log.Fatal(err)
	}

	// Register all services
	// User service
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	if err := auth_supertokens.Init(userService); err != nil {
		log.Fatal("Failed to initialize supertokens", err)
	}
	// Account service
	accountRepo := account.NewRepository(db)
	accountService := account.NewService(accountRepo)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors.Config())
	router.Use(auth.Verify())

	// Register all routes
	apiRoute := router.Group("api")
	user.SetHandlers(apiRoute, userService)
	account.SetHandlers(apiRoute, accountService)

	// Misc endpoint
	apiRoute.GET("/sessioninfo", auth.VerifySession(nil), auth.SessionInfo)
	apiRoute.POST("/account-migrate", database.ValidateKey(), database.MigrationHandler)
	apiRoute.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// starting the server
	addr := fmt.Sprintf(":%s", os.Getenv("APP_PORT"))
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start account service: %v", err)
	}
}

package database

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tryoasnafi/be-assignment/account/internal/user"
	"github.com/tryoasnafi/be-assignment/common/model"
	"gorm.io/gorm"
)

func Migrate(d *gorm.DB) error {
	return d.AutoMigrate(
		user.User{},
		model.Account{},
		model.AccountHistories{},
	)
}

type MigrationKey struct {
	Key string `json:"key" example:"helloworld123"`
}

type DefaultResponse struct {
	Message string `json:"message"`
}

// Validatekey is middleware to validate the key of migration
func ValidateKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		in := MigrationKey{}
		key := os.Getenv("APP_MIGRATE_KEY")
		if key == "" {
			key = "helloworld"
		}
		if err := c.ShouldBindJSON(&in); err != nil || in.Key != key {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid key",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Migration handler for account service
// @Summary migrate account schema
// @Schemes
// @Description migrate account schema and the related tables
// @Tags migration
// @Accept json
// @Produce json
// @Success 200 {object} DefaultResponse
// @Router /account-migrate [post]
// @Param request body MigrationKey true "key"
func MigrationHandler(c *gin.Context) {
	if err := Migrate(DB); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Account Service Migration Success",
	})
}

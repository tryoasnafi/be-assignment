package database

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tryoasnafi/be-assignment/common/dto"
	"gorm.io/gorm"
)

func Migrate(d *gorm.DB) error {
	return d.AutoMigrate(
		dto.Transaction{},
	)
}

// Validatekey is middleware to validate the key of migration
func ValidateKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		in := struct {
			Key string `json:"key"`
		}{}
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

// Handler for exec migration
func MigrationHandler(c *gin.Context) {
	if err := Migrate(DB); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Transaction Service Migration Success",
	})
}

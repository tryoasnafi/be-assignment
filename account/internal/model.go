package internal

import (
	"time"

	"gorm.io/gorm"
)

// Copied from gorm.Model with json tag modification
type Model struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

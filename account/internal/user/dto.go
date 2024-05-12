package user

import (
	"github.com/google/uuid"
	"github.com/tryoasnafi/be-assignment/account/internal"
	"github.com/tryoasnafi/be-assignment/common/dto"
	"gorm.io/datatypes"
)

type User struct {
	internal.Model
	UUID        uuid.UUID      `json:"uuid" gorm:"uniqueIndex;type:uuid;default:gen_random_uuid()"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Address     string         `json:"address"`
	DOB         datatypes.Date `json:"dob"`
	Email       string         `json:"email"`
	PhoneNumber string         `json:"phone_number"`
	Accounts    []*dto.Account `json:"accounts,omitempty" gorm:"foreignKey:user_id;references:uuid"`
}

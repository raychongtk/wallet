package user

import (
	"github.com/google/uuid"
	"time"
)

type AppUser struct {
	ID          uuid.UUID
	Email       string
	PhoneNumber string
	Password    string
	FirstName   string
	LastName    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (user AppUser) TableName() string {
	return "app_user"
}

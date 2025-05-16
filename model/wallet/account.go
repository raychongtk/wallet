package wallet

import (
	"github.com/google/uuid"
	"time"
)

type Account struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	AccountType string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (account Account) TableName() string {
	return "account"
}

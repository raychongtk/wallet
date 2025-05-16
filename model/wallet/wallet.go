package wallet

import (
	"github.com/google/uuid"
	"time"
)

type Wallet struct {
	ID           uuid.UUID
	AccountID    uuid.UUID
	Currency     string
	DecimalPlace int
	WalletStatus string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (wallet Wallet) TableName() string {
	return "wallet"
}

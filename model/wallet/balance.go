package wallet

import (
	"github.com/google/uuid"
	"time"
)

type Balance struct {
	ID          uuid.UUID
	WalletID    uuid.UUID
	BalanceType string
	Balance     int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (balance Balance) TableName() string {
	return "balance"
}

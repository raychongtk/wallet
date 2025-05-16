package wallet

import (
	"github.com/google/uuid"
	"time"
)

type BalanceType string

const (
	RESERVED_DEBIT  BalanceType = "RESERVED_DEBIT"
	RESERVED_CREDIT BalanceType = "RESERVED_CREDIT"
	COMMITTED       BalanceType = "COMMITTED"
)

type Balance struct {
	ID          uuid.UUID
	WalletID    uuid.UUID
	BalanceType BalanceType
	Balance     int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (balance Balance) TableName() string {
	return "balance"
}

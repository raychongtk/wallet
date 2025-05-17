package movement

import (
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	ID          uuid.UUID
	MovementID  uuid.UUID
	WalletID    uuid.UUID
	BalanceType string
	Balance     int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (transaction Transaction) TableName() string {
	return "transaction"
}

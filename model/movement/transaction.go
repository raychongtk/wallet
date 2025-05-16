package movement

import (
	"github.com/google/uuid"
	"github.com/raychongtk/wallet/model/wallet"
	"time"
)

type Transaction struct {
	ID          uuid.UUID
	MovementID  uuid.UUID
	WalletID    uuid.UUID
	BalanceType wallet.BalanceType
	Balance     int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (transaction Transaction) TableName() string {
	return "transaction"
}

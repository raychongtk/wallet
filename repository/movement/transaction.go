package movement

import (
	"github.com/google/uuid"
	"github.com/raychongtk/wallet/repository/wallet"
	"time"
)

type Transaction struct {
	ID             uuid.UUID
	DebitWalletID  uuid.UUID
	CreditWalletID uuid.UUID
	BalanceType    wallet.BalanceType
	Balance        int
	MovementStatus string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (transaction Transaction) TableName() string {
	return "transaction"
}

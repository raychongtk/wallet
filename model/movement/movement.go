package movement

import (
	"github.com/google/uuid"
	"time"
)

type Movement struct {
	ID             uuid.UUID
	DebitWalletID  uuid.UUID
	CreditWalletID uuid.UUID
	DebitBalance   int
	CreditBalance  int
	MovementStatus string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (movement Movement) TableName() string {
	return "movement"
}

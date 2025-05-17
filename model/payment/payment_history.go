package payment

import (
	"github.com/google/uuid"
	"time"
)

type PaymentHistory struct {
	ID          uuid.UUID
	PayerUserId string
	PayerName   string
	PayeeUserId string
	PayeeName   string
	Amount      int
	PayType     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (paymentHistory PaymentHistory) TableName() string {
	return "payment_history"
}

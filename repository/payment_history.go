package repository

import (
	"github.com/raychongtk/wallet/model/payment"
	"gorm.io/gorm"
)

type PaymentHistoryRepository interface {
	CreatePaymentHistory(db *gorm.DB, paymentHistory *payment.PaymentHistory) (*payment.PaymentHistory, error)
}

type PgPaymentHistoryRepository struct {
	db *gorm.DB
}

func ProvidePaymentHistoryRepository(db gorm.DB) PaymentHistoryRepository {
	return &PgPaymentHistoryRepository{&db}
}

func (m *PgPaymentHistoryRepository) CreatePaymentHistory(db *gorm.DB, paymentHistory *payment.PaymentHistory) (*payment.PaymentHistory, error) {
	result := db.Create(paymentHistory)
	if result.Error != nil {
		return nil, result.Error
	}
	return paymentHistory, nil
}

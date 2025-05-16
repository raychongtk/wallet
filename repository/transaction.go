package repository

import (
	"github.com/raychongtk/wallet/model/movement"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransactions(transaction []movement.Transaction) error
}

type PgTransactionRepository struct {
	db *gorm.DB
}

func ProvideTransactionRepository(db gorm.DB) TransactionRepository {
	return &PgTransactionRepository{&db}
}

func (m *PgTransactionRepository) CreateTransactions(transactions []movement.Transaction) error {
	result := m.db.Create(transactions)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

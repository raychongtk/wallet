package repository

import (
	"github.com/google/uuid"
	"github.com/raychongtk/wallet/model/wallet"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BalanceRepository interface {
	UpdateBalance(walletID uuid.UUID, balance int, balanceType wallet.BalanceType) error
}

type PgBalanceRepository struct {
	db *gorm.DB
}

func ProvideBalanceRepository(db gorm.DB) BalanceRepository {
	return &PgBalanceRepository{&db}
}

func (m *PgBalanceRepository) UpdateBalance(walletID uuid.UUID, balance int, balanceType wallet.BalanceType) error {
	tx := m.db.Begin()
	walletBalance, err := m.GetBalance(walletID, balanceType)
	if err != nil {
		return err
	}
	walletBalance.Balance += balance
	result := m.db.Save(walletBalance)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	tx.Commit()
	return nil
}

func (m *PgBalanceRepository) GetBalance(walletID uuid.UUID, balanceType wallet.BalanceType) (*wallet.Balance, error) {
	var balance wallet.Balance
	result := m.db.Where("wallet_id = ? AND balance_type = ?", walletID.String(), balanceType).Clauses(clause.Locking{Strength: "FOR UPDATE"}).First(&balance)
	if result.Error != nil {
		return nil, result.Error
	}
	return &balance, nil
}

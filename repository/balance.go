package repository

import (
	"errors"
	"github.com/google/uuid"
	"github.com/raychongtk/wallet/model/wallet"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BalanceRepository interface {
	AddBalance(db *gorm.DB, walletID uuid.UUID, balance int, balanceType string) error
	DeductBalance(db *gorm.DB, walletID uuid.UUID, balance int, balanceType string, accountType string) error
	GetBalanceWithLock(db *gorm.DB, walletID uuid.UUID, balanceType string) (*wallet.Balance, error)
	GetBalance(walletID uuid.UUID, balanceType string) (*wallet.Balance, error)
}

type PgBalanceRepository struct {
	db *gorm.DB
}

func ProvideBalanceRepository(db gorm.DB) BalanceRepository {
	return &PgBalanceRepository{&db}
}

func (m *PgBalanceRepository) AddBalance(db *gorm.DB, walletID uuid.UUID, balance int, balanceType string) error {
	walletBalance, err := m.GetBalanceWithLock(db, walletID, balanceType)
	if err != nil {
		return err
	}
	walletBalance.Balance += balance
	result := db.Save(walletBalance)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *PgBalanceRepository) DeductBalance(db *gorm.DB, walletID uuid.UUID, balance int, balanceType string, accountType string) error {
	walletBalance, err := m.GetBalanceWithLock(db, walletID, balanceType)
	if err != nil {
		return err
	}
	if accountType == "CUSTOMER" && walletBalance.Balance < balance {
		return errors.New("insufficient balance")
	}
	walletBalance.Balance -= balance
	result := db.Save(walletBalance)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *PgBalanceRepository) GetBalanceWithLock(db *gorm.DB, walletID uuid.UUID, balanceType string) (*wallet.Balance, error) {
	var balance wallet.Balance
	result := db.Clauses(clause.Locking{Strength: clause.LockingStrengthUpdate}).First(&balance, "wallet_id = ? AND balance_type = ?", walletID.String(), balanceType)
	if result.Error != nil {
		return nil, result.Error
	}
	return &balance, nil
}

func (m *PgBalanceRepository) GetBalance(walletID uuid.UUID, balanceType string) (*wallet.Balance, error) {
	var balance wallet.Balance
	result := m.db.First(&balance, "wallet_id = ? AND balance_type = ?", walletID.String(), balanceType)
	if result.Error != nil {
		return nil, result.Error
	}
	return &balance, nil
}

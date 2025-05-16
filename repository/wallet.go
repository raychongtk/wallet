package repository

import (
	"errors"
	"github.com/google/uuid"
	"github.com/raychongtk/wallet/model/wallet"
	"gorm.io/gorm"
)

type WalletRepository interface {
	GetWallet(accountId uuid.UUID) (*wallet.Wallet, error)
}

type PgWalletRepository struct {
	db *gorm.DB
}

func ProvideWalletRepository(db gorm.DB) WalletRepository {
	return &PgWalletRepository{&db}
}

func (m *PgWalletRepository) GetWallet(accountId uuid.UUID) (*wallet.Wallet, error) {
	appUser, err := m.find(accountId)
	if err != nil {
		return nil, err
	}
	return appUser, nil
}

func (m *PgWalletRepository) find(accountId uuid.UUID) (*wallet.Wallet, error) {
	var appWallet wallet.Wallet
	result := m.db.Where("account_id = ?", accountId.String()).Find(&appWallet)
	if result.Error != nil {
		panic(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("account not found")
	}

	return &appWallet, nil
}

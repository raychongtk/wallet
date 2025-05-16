package repository

import (
	"errors"
	"github.com/google/uuid"
	"github.com/raychongtk/wallet/model/wallet"
	"gorm.io/gorm"
)

type AccountRepository interface {
	GetAccount(userId uuid.UUID) (*wallet.Account, error)
}

type PgAccountRepository struct {
	db *gorm.DB
}

func ProvideAccountRepository(db gorm.DB) AccountRepository {
	return &PgAccountRepository{&db}
}

func (m *PgAccountRepository) GetAccount(userId uuid.UUID) (*wallet.Account, error) {
	appUser, err := m.find(userId)
	if err != nil {
		return nil, err
	}
	return appUser, nil
}

func (m *PgAccountRepository) find(userId uuid.UUID) (*wallet.Account, error) {
	var account wallet.Account
	result := m.db.Where("user_id = ?", userId.String()).Find(&account)
	if result.Error != nil {
		panic(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("account not found")
	}

	return &account, nil
}

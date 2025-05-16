package repository

import (
	"errors"
	"github.com/google/uuid"
	"github.com/raychongtk/wallet/model/user"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUser(id uuid.UUID) (*user.AppUser, error)
}

type PgUserRepository struct {
	db *gorm.DB
}

func ProvideUserRepository(db gorm.DB) UserRepository {
	return &PgUserRepository{&db}
}

func (m *PgUserRepository) GetUser(id uuid.UUID) (*user.AppUser, error) {
	appUser, err := m.find(id)
	if err != nil {
		return nil, err
	}
	return appUser, nil
}

func (m *PgUserRepository) find(id uuid.UUID) (*user.AppUser, error) {
	var appUser user.AppUser
	result := m.db.Where("id = ?", id.String()).Find(&appUser)
	if result.Error != nil {
		panic(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}

	return &appUser, nil
}

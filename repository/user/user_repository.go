package user

import (
	"errors"
	"github.com/google/uuid"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var (
	WireSet = wire.NewSet(NewRepository)
)

type UserRepository interface {
	GetAccount(id uuid.UUID) (*AppUser, error)
}

type PgUserRepository struct {
	db *gorm.DB
}

func NewRepository(db gorm.DB) UserRepository {
	return &PgUserRepository{&db}
}

func (m *PgUserRepository) GetAccount(id uuid.UUID) (*AppUser, error) {
	user, err := m.find(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m *PgUserRepository) find(id uuid.UUID) (*AppUser, error) {
	var appUser AppUser
	result := m.db.Where("id = ?", id.String()).Find(&appUser)
	if result.Error != nil {
		panic(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}

	return &appUser, nil
}

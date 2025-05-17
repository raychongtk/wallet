package repository

import (
	"github.com/raychongtk/wallet/model/movement"
	"gorm.io/gorm"
)

type MovementRepository interface {
	CreateMovement(db *gorm.DB, movement *movement.Movement) (*movement.Movement, error)
}

type PgMovementRepository struct {
	db *gorm.DB
}

func ProvideMovementRepository(db gorm.DB) MovementRepository {
	return &PgMovementRepository{&db}
}

func (m *PgMovementRepository) CreateMovement(db *gorm.DB, movement *movement.Movement) (*movement.Movement, error) {
	result := db.Create(movement)
	if result.Error != nil {
		return nil, result.Error
	}
	return movement, nil
}

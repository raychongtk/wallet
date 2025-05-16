package repository

import (
	"errors"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/raychongtk/go-web/util"
	"gorm.io/gorm"
)

var (
	WireSet = wire.NewSet(NewRepository)
)

type UserRepository interface {
	Authenticate(username string, password string) (bool, error)
	CreateAccount(username string, password string, firstName string, lastName string) (bool, error)
	GetAccount(username string) (*AppUser, error)
}

type PgUserRepository struct {
	db *gorm.DB
}

func NewRepository(db gorm.DB) UserRepository {
	return &PgUserRepository{&db}
}

func (m *PgUserRepository) GetAccount(username string) (*AppUser, error) {
	user, err := m.find(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m *PgUserRepository) Authenticate(username string, password string) (bool, error) {
	user, err := m.find(username)
	if err != nil {
		return false, err
	}
	return util.ValidateHash(user.Password, password), nil
}

func (m *PgUserRepository) CreateAccount(username string, password string, firstName string, lastName string) (bool, error) {
	exists := m.exists(username)
	if exists {
		return false, errors.New("username already in use")
	}

	user := &AppUser{
		ID:        uuid.New().String(),
		Username:  username,
		Password:  util.Hash(password),
		FirstName: firstName,
		LastName:  lastName,
	}
	result := m.db.Create(user)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (m *PgUserRepository) find(username string) (*AppUser, error) {
	var appUser AppUser
	result := m.db.Where("username = ?", username).Find(&appUser)
	if result.Error != nil {
		panic(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}

	return &appUser, nil
}

func (m *PgUserRepository) exists(username string) bool {
	var count int64
	m.db.Table("app_user").Where("username = ?", username).Count(&count)
	return count > 0
}

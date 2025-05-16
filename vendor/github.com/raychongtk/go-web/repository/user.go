package repository

import "time"

type AppUser struct {
	ID        string
	Username  string
	Password  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (user AppUser) TableName() string {
	return "app_user"
}

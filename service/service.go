package service

import (
	"github.com/google/wire"
	"github.com/raychongtk/wallet/repository"
)

var (
	WireSet = wire.NewSet(ProvideService)
)

type Service struct {
	userRepo repository.UserRepository
}

func ProvideService(repository repository.UserRepository) (*Service, error) {
	return &Service{userRepo: repository}, nil
}

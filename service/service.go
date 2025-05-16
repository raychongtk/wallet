package service

import (
	"github.com/google/wire"
	"github.com/raychongtk/wallet/repository"
	"gorm.io/gorm"
)

var (
	WireSet = wire.NewSet(ProvideService)
)

type Service struct {
	userRepo        repository.UserRepository
	movementRepo    repository.MovementRepository
	accountRepo     repository.AccountRepository
	walletRepo      repository.WalletRepository
	transactionRepo repository.TransactionRepository
	balanceRepo     repository.BalanceRepository
	db              *gorm.DB
}

func ProvideService(
	userRepo repository.UserRepository,
	movementRepo repository.MovementRepository,
	accountRepo repository.AccountRepository,
	walletRepo repository.WalletRepository,
	transactionRepo repository.TransactionRepository,
	balanceRepo repository.BalanceRepository,
	db *gorm.DB,
) (*Service, error) {
	return &Service{
		userRepo:        userRepo,
		movementRepo:    movementRepo,
		accountRepo:     accountRepo,
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		balanceRepo:     balanceRepo,
		db:              db,
	}, nil
}

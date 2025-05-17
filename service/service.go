package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/raychongtk/wallet/repository"
	"gorm.io/gorm"
)

var (
	WireSet = wire.NewSet(ProvideService, ProvideRoutes)
)

type Service struct {
	userRepo           repository.UserRepository
	movementRepo       repository.MovementRepository
	accountRepo        repository.AccountRepository
	walletRepo         repository.WalletRepository
	transactionRepo    repository.TransactionRepository
	balanceRepo        repository.BalanceRepository
	paymentHistoryRepo repository.PaymentHistoryRepository
	db                 gorm.DB
}

func ProvideService(
	userRepo repository.UserRepository,
	movementRepo repository.MovementRepository,
	accountRepo repository.AccountRepository,
	walletRepo repository.WalletRepository,
	transactionRepo repository.TransactionRepository,
	balanceRepo repository.BalanceRepository,
	paymentHistoryRepo repository.PaymentHistoryRepository,
	db gorm.DB,
) (*Service, error) {
	return &Service{
		userRepo:           userRepo,
		movementRepo:       movementRepo,
		accountRepo:        accountRepo,
		walletRepo:         walletRepo,
		transactionRepo:    transactionRepo,
		balanceRepo:        balanceRepo,
		paymentHistoryRepo: paymentHistoryRepo,
		db:                 db,
	}, nil
}

func ProvideRoutes(service *Service) *gin.Engine {
	r := gin.New()

	r.POST("/api/v1/wallet/deposit", service.Deposit)
	r.POST("/api/v1/wallet/withdrawal", service.Withdraw)
	r.POST("/api/v1/wallet/transfer", service.Transfer)
	r.GET("/api/v1/wallet/balance", service.GetBalance)
	r.GET("/api/v1/wallet/payment-history", service.GetPaymentHistory)
	return r
}

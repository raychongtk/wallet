package service

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
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
	memoryStore        redis.Client
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
	memoryStore redis.Client,
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
		memoryStore:        memoryStore,
	}, nil
}

func ProvideRoutes(service *Service) *gin.Engine {
	r := gin.New()

	protected := r.Group("/api/v1/wallet")
	protected.Use(service.ValidateRequestID())
	protected.POST("/deposit", service.Deposit)
	protected.POST("/withdrawal", service.Withdraw)
	protected.POST("/transfer", service.Transfer)

	r.GET("/api/v1/wallet/balance", service.GetBalance)
	r.GET("/api/v1/wallet/payment-history", service.GetPaymentHistory)
	return r
}

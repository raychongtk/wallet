package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/raychongtk/wallet/model/movement"
	"github.com/raychongtk/wallet/model/wallet"
	"github.com/raychongtk/wallet/util"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func (s *Service) Withdraw(ctx *gin.Context) {
	var req WithdrawalRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		util.Error("Invalid params", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, &WithdrawalResponse{Result: false, ErrorCode: "INVALID_PARAMETERS"})
		return
	}
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		util.Error("Invalid user", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, &WithdrawalResponse{Result: false, ErrorCode: "INVALID_ACCOUNT"})
		return
	}
	balance, err := util.ConvertToInt(req.Balance)
	if err != nil || balance <= 0 {
		ctx.JSON(http.StatusBadRequest, &TransferResponse{Result: false, ErrorCode: "INVALID_PARAMETERS"})
		return
	}

	appUser, err := s.userRepo.GetUser(userId)
	if err != nil {
		util.Error("Invalid user", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, &WithdrawalResponse{Result: false, ErrorCode: "INVALID_ACCOUNT"})
		return
	}
	account, err := s.accountRepo.GetAccount(appUser.ID)
	if err != nil {
		util.Error("Invalid account", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, &WithdrawalResponse{Result: false, ErrorCode: "INVALID_ACCOUNT"})
		return
	}
	userWallet, err := s.walletRepo.GetWallet(account.ID)
	if err != nil {
		util.Error("Invalid wallet", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, &WithdrawalResponse{Result: false, ErrorCode: "INVALID_ACCOUNT"})
		return
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, &WithdrawalResponse{Result: false, ErrorCode: "INTERNAL_ERROR"})
		}
	}()
	balance = balance * 100
	groupId := uuid.New()
	newMovement := &movement.Movement{
		ID:             uuid.New(),
		GroupID:        groupId,
		DebitWalletID:  util.GetLiabilityAccount(),
		CreditWalletID: userWallet.ID,
		DebitBalance:   balance,
		CreditBalance:  -balance,
		MovementStatus: "COMPLETED",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	createdMovement, err := s.movementRepo.CreateMovement(tx, newMovement)
	if err != nil {
		tx.Rollback()
		util.Error("Create movement failed", zap.String("user_id", userId.String()), zap.Int("balance", balance), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, &WithdrawalResponse{Result: false, ErrorCode: "INTERNAL_ERROR"})
		return
	}
	transactions := GenerateWithdrawalTransactions(balance, userWallet.ID, createdMovement.ID)
	err = s.transactionRepo.CreateTransactions(tx, transactions)
	if err != nil {
		tx.Rollback()
		util.Error("Create payment transaction failed", zap.String("user_id", userId.String()), zap.Int("balance", balance), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, &WithdrawalResponse{Result: false, ErrorCode: "INTERNAL_ERROR"})
		return
	}
	committed := commitWithdrawalBalance(s, userWallet, balance, tx)

	if !committed {
		tx.Rollback()
		util.Error("Update balance failed", zap.String("user_id", userId.String()), zap.Int("balance", balance), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, &WithdrawalResponse{Result: false, ErrorCode: "INTERNAL_ERROR"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		util.Error("Commit transaction failed", zap.String("user_id", userId.String()), zap.Int("balance", balance), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, &WithdrawalResponse{Result: false, ErrorCode: "INTERNAL_ERROR"})
		return
	}

	util.Info("Withdrawal successfully", zap.String("user_id", userId.String()), zap.Int("balance", balance))
	ctx.JSON(http.StatusOK, &WithdrawalResponse{Result: true})
}

func commitWithdrawalBalance(s *Service, wallet *wallet.Wallet, balance int, tx *gorm.DB) bool {
	// Skip reserved balance because we don't need to wait for external clearing operations
	// customer asset decreased
	customerAccountErr := s.balanceRepo.DeductBalance(tx, wallet.ID, balance, "COMMITTED", "CUSTOMER")
	if customerAccountErr != nil {
		tx.Rollback()
		return false
	}
	// company liability increased
	chartAccountError := s.balanceRepo.AddBalance(tx, util.GetLiabilityAccount(), balance, "COMMITTED")
	if chartAccountError != nil {
		tx.Rollback()
		return false
	}
	return true
}

func GenerateWithdrawalTransactions(balance int, creditWalletID uuid.UUID, movementID uuid.UUID) []movement.Transaction {
	var transactions []movement.Transaction

	debitTransaction := movement.Transaction{
		ID:          uuid.New(),
		MovementID:  movementID,
		WalletID:    util.GetLiabilityAccount(),
		BalanceType: "COMMITTED",
		Balance:     balance,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	transactions = append(transactions, debitTransaction)

	creditTransaction := movement.Transaction{
		ID:          uuid.New(),
		WalletID:    creditWalletID,
		BalanceType: "COMMITTED",
		Balance:     -balance,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	transactions = append(transactions, creditTransaction)

	return transactions
}

type WithdrawalRequest struct {
	UserId  string `json:"user_id" binding:"required"`
	Balance string `json:"balance" binding:"required"`
}

type WithdrawalResponse struct {
	Result    bool   `json:"result" binding:"required"`
	ErrorCode string `json:"error_code,omitempty"`
}

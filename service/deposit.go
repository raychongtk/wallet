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

func (s *Service) Deposit(ctx *gin.Context) {
	var req depositRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		util.Error("Invalid params", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, &depositResponse{Result: false, ErrorCode: "INVALID_PARAMETERS"})
		return
	}
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		util.Error("Invalid user", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, &depositResponse{Result: false, ErrorCode: "INVALID_ACCOUNT"})
		return
	}
	appUser, err := s.userRepo.GetUser(userId)
	if err != nil {
		util.Error("Invalid user", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, &depositResponse{Result: false, ErrorCode: "INVALID_ACCOUNT"})
		return
	}
	account, err := s.accountRepo.GetAccount(appUser.ID)
	if err != nil {
		util.Error("Invalid account", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, &depositResponse{Result: false, ErrorCode: "INVALID_ACCOUNT"})
		return
	}
	userWallet, err := s.walletRepo.GetWallet(account.ID)
	if err != nil {
		util.Error("Invalid wallet", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, &depositResponse{Result: false, ErrorCode: "INVALID_ACCOUNT"})
		return
	}
	balance, err := util.ConvertToInt(req.Balance)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &depositResponse{Result: false, ErrorCode: "INVALID_PARAMETERS"})
		return
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, &depositResponse{Result: false, ErrorCode: "INTERNAL_ERROR"})
		}
	}()
	balance = balance * 100
	newMovement := &movement.Movement{
		ID:             uuid.New(),
		DebitWalletID:  util.GetAssetAccount(),
		CreditWalletID: userWallet.ID,
		DebitBalance:   balance,
		CreditBalance:  balance,
		MovementStatus: "COMPLETED",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	createdMovement, err := s.movementRepo.CreateMovement(tx, newMovement)
	if err != nil {
		tx.Rollback()
		util.Error("Create movement failed", zap.String("user_id", userId.String()), zap.Int("balance", balance), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, &depositResponse{Result: false, ErrorCode: "CREATE_MOVEMENT_FAILED"})
		return
	}
	transactions := GenerateTransactions(balance, userWallet.ID, createdMovement.ID)
	err = s.transactionRepo.CreateTransactions(tx, transactions)
	if err != nil {
		tx.Rollback()
		util.Error("Create payment transaction failed", zap.String("user_id", userId.String()), zap.Int("balance", balance), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, &depositResponse{Result: false, ErrorCode: "CREATE_TRANSACTION_FAILED"})
		return
	}
	committed := commitBalance(s, userWallet, balance, tx)

	if !committed {
		tx.Rollback()
		util.Error("Update balance failed", zap.String("user_id", userId.String()), zap.Int("balance", balance), zap.Error(err))
		ctx.JSON(http.StatusBadRequest, &depositResponse{Result: false, ErrorCode: "UPDATE_BALANCE_FAILED"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		util.Error("Commit transaction failed", zap.String("user_id", userId.String()), zap.Int("balance", balance), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, &depositResponse{Result: false, ErrorCode: "COMMIT_FAILED"})
		return
	}

	util.Info("Deposit successfully", zap.String("user_id", userId.String()), zap.Int("balance", balance))
	ctx.JSON(http.StatusOK, &depositResponse{Result: true})
}

func commitBalance(s *Service, wallet *wallet.Wallet, balance int, tx *gorm.DB) bool {
	// Skip reserved balance because we don't need to wait for external clearing operations
	customerAccountErr := s.balanceRepo.UpdateBalance(tx, wallet.ID, balance, "COMMITTED")
	if customerAccountErr != nil {
		tx.Rollback()
		return false
	}
	chartAccountError := s.balanceRepo.UpdateBalance(tx, util.GetAssetAccount(), balance, "COMMITTED")
	if chartAccountError != nil {
		tx.Rollback()
		return false
	}
	return true
}

func GenerateTransactions(balance int, creditWalletID uuid.UUID, movementID uuid.UUID) []movement.Transaction {
	var transactions []movement.Transaction

	debitTransaction := movement.Transaction{
		ID:          uuid.New(),
		MovementID:  movementID,
		WalletID:    util.GetAssetAccount(),
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
		Balance:     balance,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	transactions = append(transactions, creditTransaction)

	return transactions
}

type depositRequest struct {
	UserId  string `json:"user_id" binding:"required"`
	Balance string `json:"balance" binding:"required"`
}

type depositResponse struct {
	Result    bool   `json:"result" binding:"required"`
	ErrorCode string `json:"error_code,omitempty"`
}

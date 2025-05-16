package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/raychongtk/wallet/model/movement"
	"github.com/raychongtk/wallet/util"
	"net/http"
	"time"
)

func (s *Service) Deposit(ctx *gin.Context) {
	var req depositRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &depositResponse{Result: false, ErrorCode: "INVALID_PARAMETERS"})
		return
	}
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &depositResponse{Result: false, ErrorCode: "INVALID_ACCOUNT"})
		return
	}
	appUser, err := s.userRepo.GetUser(userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &depositResponse{Result: false, ErrorCode: "INVALID_ACCOUNT"})
		return
	}
	account, err := s.accountRepo.GetAccount(appUser.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &depositResponse{Result: false, ErrorCode: "INVALID_ACCOUNT"})
		return
	}
	wallet, err := s.walletRepo.GetWallet(account.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &depositResponse{Result: false, ErrorCode: "INVALID_ACCOUNT"})
		return
	}
	balance, err := util.ConvertToInt(req.Balance)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &depositResponse{Result: false, ErrorCode: "INVALID_PARAMETERS"})
		return
	}

	// Start a database transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, &depositResponse{Result: false, ErrorCode: "INTERNAL_ERROR"})
		}
	}()

	newMovement := &movement.Movement{
		ID:             uuid.New(),
		DebitWalletID:  util.GetAssetAccount(),
		CreditWalletID: wallet.ID,
		DebitBalance:   balance,
		CreditBalance:  balance,
		MovementStatus: "COMPLETED",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	createdMovement, err := s.movementRepo.CreateMovement(newMovement)
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, &depositResponse{Result: false, ErrorCode: "CREATE_MOVEMENT_FAILED"})
		return
	}
	transactions := GenerateTransactions(balance, wallet.ID, createdMovement.ID)
	err = s.transactionRepo.CreateTransactions(transactions)
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, &depositResponse{Result: false, ErrorCode: "CREATE_TRANSACTION_FAILED"})
		return
	}
	// Update the wallet balance
	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, &depositResponse{Result: false, ErrorCode: "COMMIT_FAILED"})
		return
	}
	ctx.JSON(http.StatusOK, &depositResponse{Result: true})
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

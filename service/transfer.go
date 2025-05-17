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

func (s *Service) Transfer(ctx *gin.Context) {
	var req TransferRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		util.Error("Invalid params", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, &TransferResponse{Result: false, ErrorCode: "INVALID_PARAMETERS"})
		return
	}
	if req.CreditUserId == req.DebitUserId {
		ctx.JSON(http.StatusBadRequest, &TransferResponse{Result: false, ErrorCode: "CANNOT_TRANSFER_TO_SELF"})
		return
	}
	balance, err := util.ConvertToInt(req.Balance)
	if err != nil || balance <= 0 {
		ctx.JSON(http.StatusBadRequest, &TransferResponse{Result: false, ErrorCode: "INVALID_PARAMETERS"})
		return
	}

	creditUserId, isValid := validUserId(req.CreditUserId)
	if !isValid {
		ctx.JSON(http.StatusBadRequest, &TransferResponse{Result: false, ErrorCode: "INVALID_ACCOUNT"})
		return
	}
	debitUserId, isValid := validUserId(req.DebitUserId)
	if !isValid {
		ctx.JSON(http.StatusBadRequest, &TransferResponse{Result: false, ErrorCode: "INVALID_ACCOUNT"})
		return
	}
	util.With(zap.String("credit_user_id", creditUserId.String()), zap.String("debit_user_id", debitUserId.String()))

	creditAppUser, err := s.userRepo.GetUser(creditUserId)
	if err != nil {
		util.Error("Invalid credit user", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, &TransferResponse{Result: false, ErrorCode: "INVALID_ACCOUNT"})
		return
	}
	debitAppUser, err := s.userRepo.GetUser(debitUserId)
	if err != nil {
		util.Error("Invalid debit user", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, &TransferResponse{Result: false, ErrorCode: "INVALID_ACCOUNT"})
		return
	}

	creditAccount, err := s.accountRepo.GetAccount(creditAppUser.ID)
	if err != nil {
		util.Error("Invalid credit account", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, &TransferResponse{Result: false, ErrorCode: "INVALID_ACCOUNT"})
		return
	}

	debitAccount, err := s.accountRepo.GetAccount(debitAppUser.ID)
	if err != nil {
		util.Error("Invalid debit account", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, &TransferResponse{Result: false, ErrorCode: "INVALID_ACCOUNT"})
		return
	}

	creditUserWallet, err := s.walletRepo.GetWallet(creditAccount.ID)
	if err != nil {
		util.Error("Invalid credit wallet", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, &TransferResponse{Result: false, ErrorCode: "INVALID_ACCOUNT"})
		return
	}
	debitUserWallet, err := s.walletRepo.GetWallet(debitAccount.ID)
	if err != nil {
		util.Error("Invalid debit wallet", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, &TransferResponse{Result: false, ErrorCode: "INVALID_ACCOUNT"})
		return
	}

	util.With(zap.String("balance", req.Balance))
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, &TransferResponse{Result: false, ErrorCode: "INTERNAL_ERROR"})
		}
	}()
	groupId := uuid.New()
	var movements []movement.Movement
	creditMovement := movement.Movement{
		ID:             uuid.New(),
		GroupID:        groupId,
		DebitWalletID:  util.GetLiabilityAccount(),
		CreditWalletID: creditUserWallet.ID,
		DebitBalance:   balance,
		CreditBalance:  -balance,
		MovementStatus: "COMPLETED",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	debitMovement := movement.Movement{
		ID:             uuid.New(),
		GroupID:        groupId,
		DebitWalletID:  debitUserWallet.ID,
		CreditWalletID: util.GetLiabilityAccount(),
		DebitBalance:   balance,
		CreditBalance:  -balance,
		MovementStatus: "COMPLETED",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	movements = append(movements, debitMovement)
	movements = append(movements, creditMovement)

	err = s.movementRepo.CreateMovements(tx, movements)
	if err != nil {
		tx.Rollback()
		util.Error("Create movement failed")
		ctx.JSON(http.StatusBadRequest, &TransferResponse{Result: false, ErrorCode: "INTERNAL_ERROR"})
		return
	}
	var transferTransactions []movement.Transaction
	transferOurTransactions := GenerateTransferOutTransactions(balance, creditUserWallet.ID, creditMovement.ID)
	transferTransactions = append(transferTransactions, transferOurTransactions...)
	transferInTransactions := GenerateTransferInTransactions(balance, debitUserWallet.ID, debitMovement.ID)
	transferTransactions = append(transferTransactions, transferInTransactions...)
	err = s.transactionRepo.CreateTransactions(tx, transferTransactions)
	if err != nil {
		tx.Rollback()
		util.Error("Create movement failed")
		ctx.JSON(http.StatusBadRequest, &TransferResponse{Result: false, ErrorCode: "INTERNAL_ERROR"})
		return
	}
	committed := commitTransferBalance(s, debitUserWallet, creditUserWallet, balance, tx)

	if !committed {
		tx.Rollback()
		util.Error("Create movement failed")
		ctx.JSON(http.StatusBadRequest, &TransferResponse{Result: false, ErrorCode: "INTERNAL_ERROR"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		util.Error("Create movement failed")
		ctx.JSON(http.StatusInternalServerError, &TransferResponse{Result: false, ErrorCode: "INTERNAL_ERROR"})
		return
	}

	util.Info("Transfer successfully")
	ctx.JSON(http.StatusOK, &TransferResponse{Result: true})
}

func validUserId(userId string) (uuid.UUID, bool) {
	creditUserId, err := uuid.Parse(userId)
	if err != nil {
		util.Error("Invalid user", zap.Error(err))
		return uuid.UUID{}, false
	}
	return creditUserId, true
}

func commitTransferBalance(s *Service, debitWallet *wallet.Wallet, creditWallet *wallet.Wallet, balance int, tx *gorm.DB) bool {
	// Skip reserved balance because we don't need to wait for external clearing operations
	// customer asset decreased
	creditCustomerAccountErr := s.balanceRepo.DeductBalance(tx, creditWallet.ID, balance, "COMMITTED", "CUSTOMER")
	if creditCustomerAccountErr != nil {
		tx.Rollback()
		return false
	}
	// company liability increased
	debitChartAccountError := s.balanceRepo.AddBalance(tx, util.GetLiabilityAccount(), balance, "COMMITTED")
	if debitChartAccountError != nil {
		tx.Rollback()
		return false
	}

	// customer asset increased
	debitCustomerAccountErr := s.balanceRepo.AddBalance(tx, debitWallet.ID, balance, "COMMITTED")
	if debitCustomerAccountErr != nil {
		tx.Rollback()
		return false
	}
	// company liability decreased
	creditChartAccountError := s.balanceRepo.DeductBalance(tx, util.GetLiabilityAccount(), balance, "COMMITTED", "CHART")
	if creditChartAccountError != nil {
		tx.Rollback()
		return false
	}
	return true
}

func GenerateTransferOutTransactions(balance int, creditWalletID uuid.UUID, movementID uuid.UUID) []movement.Transaction {
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

func GenerateTransferInTransactions(balance int, debitWalletID uuid.UUID, movementID uuid.UUID) []movement.Transaction {
	var transactions []movement.Transaction

	debitTransaction := movement.Transaction{
		ID:          uuid.New(),
		MovementID:  movementID,
		WalletID:    debitWalletID,
		BalanceType: "COMMITTED",
		Balance:     balance,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	transactions = append(transactions, debitTransaction)

	creditTransaction := movement.Transaction{
		ID:          uuid.New(),
		WalletID:    util.GetLiabilityAccount(),
		BalanceType: "COMMITTED",
		Balance:     -balance,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	transactions = append(transactions, creditTransaction)

	return transactions
}

type TransferRequest struct {
	CreditUserId string `json:"credit_user_id" binding:"required"`
	DebitUserId  string `json:"debit_user_id" binding:"required"`
	Balance      string `json:"balance" binding:"required"`
}

type TransferResponse struct {
	Result    bool   `json:"result" binding:"required"`
	ErrorCode string `json:"error_code,omitempty"`
}

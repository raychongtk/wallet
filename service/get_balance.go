package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/raychongtk/wallet/util"
	"go.uber.org/zap"
	"net/http"
)

func (s *Service) GetBalance(ctx *gin.Context) {
	userId, err := uuid.Parse(ctx.Query("user_id"))
	if err != nil {
		util.Error("Invalid user", zap.Error(err))
		ctx.Status(http.StatusBadRequest)
		return
	}

	appUser, err := s.userRepo.GetUser(userId)
	if err != nil {
		util.Error("Invalid user", zap.Error(err))
		ctx.Status(http.StatusBadRequest)
		return
	}
	account, err := s.accountRepo.GetAccount(appUser.ID)
	if err != nil {
		util.Error("Invalid account", zap.Error(err))
		ctx.Status(http.StatusBadRequest)
		return
	}
	userWallet, err := s.walletRepo.GetWallet(account.ID)
	if err != nil {
		util.Error("Invalid wallet", zap.Error(err))
		ctx.Status(http.StatusBadRequest)
		return
	}
	balance, err := s.balanceRepo.GetBalance(userWallet.ID, "COMMITTED")
	displayedBalance := fmt.Sprintf("%.2f", float64(balance.Balance)/100)
	ctx.JSON(http.StatusOK, &GetCustomerBalanceResponse{CustomerID: userId.String(), Currency: userWallet.Currency, Balance: displayedBalance})
}

type GetCustomerBalanceResponse struct {
	CustomerID string `json:"customer_id" binding:"required"`
	Currency   string `json:"currency" binding:"required"`
	Balance    string `json:"balance" binding:"required"`
}

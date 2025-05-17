package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/raychongtk/wallet/util"
	"go.uber.org/zap"
	"net/http"
)

func (s *Service) GetPaymentHistory(ctx *gin.Context) {
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
	histories, err := s.paymentHistoryRepo.SearchPaymentHistory(appUser.ID.String())
	if err != nil {
		util.Error("search payment history failed", zap.Error(err))
		ctx.Status(http.StatusBadRequest)
		return
	}

	var paymentHistories []PaymentHistory
	for i := 0; i < len(histories); i++ {
		paymentHistory := PaymentHistory{
			PayerName: histories[i].PayerName,
			PayeeName: histories[i].PayeeName,
			PayType:   histories[i].PayType,
			Amount:    adjustBalanceByPaymentDirection(histories[i].PayType, histories[i].PayerUserId, appUser.ID.String(), histories[i].Amount),
		}
		paymentHistories = append(paymentHistories, paymentHistory)
	}
	ctx.JSON(http.StatusOK, &SearchPaymentHistoryResponse{Histories: paymentHistories})
}

func adjustBalanceByPaymentDirection(payType string, payerUserId string, requestedUserId string, amount int) string {
	var multiplier = 1
	if payType == "TRANSFER" && payerUserId == requestedUserId {
		multiplier = -1
	} else if payType == "WITHDRAWAL" {
		multiplier = -1
	}

	return fmt.Sprintf("%.2f", float64(amount*multiplier)/100)
}

type SearchPaymentHistoryResponse struct {
	Histories []PaymentHistory `json:"histories"`
}

type PaymentHistory struct {
	PayerName string `json:"payer_name" binding:"required"`
	PayeeName string `json:"payee_name" binding:"required"`
	PayType   string `json:"pay_type" binding:"required"`
	Amount    string `json:"amount" binding:"required"`
}

package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/raychongtk/wallet/util"
	"go.uber.org/zap"
	"math"
	"net/http"
	"strconv"
)

func (s *Service) GetPaymentHistory(ctx *gin.Context) {
	userId, err := uuid.Parse(ctx.Query("user_id"))
	if err != nil {
		util.Error("Invalid user", zap.Error(err))
		ctx.Status(http.StatusBadRequest)
		return
	}
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	appUser, err := s.userRepo.GetUser(userId)
	if err != nil {
		util.Error("Invalid user", zap.Error(err))
		ctx.Status(http.StatusBadRequest)
		return
	}
	histories, err := s.paymentHistoryRepo.SearchPaymentHistory(appUser.ID.String(), page, pageSize)
	if err != nil {
		util.Error("search payment history failed", zap.Error(err))
		ctx.Status(http.StatusBadRequest)
		return
	}

	var paymentHistories []PaymentHistory
	for i := 0; i < len(histories.Data); i++ {
		paymentHistory := PaymentHistory{
			PayerName: histories.Data[i].PayerName,
			PayeeName: histories.Data[i].PayeeName,
			PayType:   histories.Data[i].PayType,
			Amount:    adjustBalanceByPaymentDirection(histories.Data[i].PayType, histories.Data[i].PayerUserId, appUser.ID.String(), histories.Data[i].Amount),
		}
		paymentHistories = append(paymentHistories, paymentHistory)
	}
	ctx.JSON(http.StatusOK, &SearchPaymentHistoryResponse{
		Histories:  paymentHistories,
		TotalCount: histories.TotalCount,
		TotalPages: calculateTotalPages(histories.TotalCount, histories.PageSize),
		Page:       histories.Page,
		PageSize:   histories.PageSize})
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

func calculateTotalPages(totalCount int64, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	return int(math.Ceil(float64(totalCount) / float64(pageSize)))
}

type SearchPaymentHistoryResponse struct {
	Histories  []PaymentHistory `json:"histories"`
	TotalCount int64            `json:"total_count"`
	TotalPages int              `json:"total_pages"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
}

type PaymentHistory struct {
	PayerName string `json:"payer_name" binding:"required"`
	PayeeName string `json:"payee_name" binding:"required"`
	PayType   string `json:"pay_type" binding:"required"`
	Amount    string `json:"amount" binding:"required"`
}

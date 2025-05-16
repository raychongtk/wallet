package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Service) Deposit(ctx *gin.Context) {
	var req depositRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid params")
		return
	}
	//s.userRepo.GetAccount()
	//s.userRepo.GetAccount()
	//created, err := s.userRepo.GetAccount(req.UserId, req.Balance)
	//if err != nil {
	//	ctx.JSON(http.StatusBadRequest, &depositResponse{Result: false, ErrorCode: "DEPOSIT_FAILED"})
	//	return
	//}
	//ctx.JSON(http.StatusOK, &depositResponse{Result: created})
}

type depositRequest struct {
	UserId  string `json:"user_id" binding:"required"`
	Balance string `json:"balance" binding:"required"`
}

type depositResponse struct {
	Result    bool   `json:"result" binding:"required"`
	ErrorCode string `json:"error_code,omitempty"`
}

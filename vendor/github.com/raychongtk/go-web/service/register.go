package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Service) Register(ctx *gin.Context) {
	var req registerAccountRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid params")
		return
	}
	created, err := s.userRepo.CreateAccount(req.Username, req.Password, req.FirstName, req.LastName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &registerAccountResponse{Result: false, ErrorCode: "CREATE_USER_FAILED"})
		return
	}
	ctx.JSON(http.StatusOK, &registerAccountResponse{Result: created})
}

type registerAccountRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type registerAccountResponse struct {
	Result    bool   `json:"result" binding:"required"`
	ErrorCode string `json:"error_code,omitempty"`
}

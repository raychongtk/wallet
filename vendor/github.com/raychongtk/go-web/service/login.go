package service

import (
	"github.com/gin-gonic/gin"
	"github.com/raychongtk/go-web/util"
	"net/http"
)

func (s *Service) Login(ctx *gin.Context) {
	var req loginRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid params")
		return
	}

	authenticated, err := s.userRepo.Authenticate(req.Username, req.Password)
	if err != nil || !authenticated {
		ctx.JSON(http.StatusUnauthorized, &loginResponse{Result: false, ErrorCode: "LOGIN_FAILED"})
		return
	}
	util.CreateSession(ctx, req.Username)
	ctx.JSON(http.StatusOK, &loginResponse{Result: authenticated})
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Result    bool   `json:"result" binding:"required"`
	ErrorCode string `json:"error_code,omitempty"`
}

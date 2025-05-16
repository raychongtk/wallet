package service

import (
	"github.com/gin-gonic/gin"
	sessions "github.com/raychongtk/go-web/util"
	"net/http"
)

func (s *Service) GetCurrentUser(ctx *gin.Context) {
	account, err := s.userRepo.GetAccount(sessions.GetUsername(ctx))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "internal error")
		return
	}
	ctx.JSON(200, &getUserResponse{Username: account.Username, FirstName: account.FirstName, LastName: account.LastName})
}

type getUserResponse struct {
	Username  string `json:"username" binding:"required"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}

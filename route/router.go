package route

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/raychongtk/wallet/service"
)

var (
	WireSet = wire.NewSet(ProvideRoutes)
)

func ProvideRoutes(service *service.Service) *gin.Engine {
	r := gin.New()

	r.POST("/api/v1/wallet/deposit", service.Deposit)
	return r
}

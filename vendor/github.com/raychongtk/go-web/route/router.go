package route

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/raychongtk/go-web/service"
)

var (
	WireSet = wire.NewSet(ProvideRoutes)
)

func ProvideRoutes(service *service.Service, store redis.Store) *gin.Engine {
	r := gin.New()
	r.Use(sessions.Sessions("session", store))

	r.PUT("/user/login", service.Login)
	r.POST("/user", service.Register)

	ajaxGroup := r.Group("/ajax")
	ajaxGroup.Use(service.VerifySession())
	ajaxGroup.GET("/user", service.GetCurrentUser)
	return r
}

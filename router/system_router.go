package router

import (
	"TPC-EDM-Server/api"
	"TPC-EDM-Server/middleware"
	"github.com/gin-gonic/gin"
)

func SysRouter(r *gin.RouterGroup) {
	app := api.App.SystemApi
	r.PUT("system/:name", middleware.AdminMiddleware, app.SystemUpdateView)
	r.GET("system", middleware.AdminMiddleware, app.SystemInfoView)
}

package router

import (
	"TPC-EDM-Server/api"
	"TPC-EDM-Server/middleware"
	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.RouterGroup) {
	app := api.App.UserApi
	r.GET("user", app.UserListView)
	r.GET("user/detail", middleware.AuthMiddleware, app.UserDetailView)
	r.POST("user", app.PwdLoginView)
	r.DELETE("user", middleware.AdminMiddleware, app.UserRemoveView)
}

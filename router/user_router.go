package router

import (
	"TPC-H-EDM-Server/api"
	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.RouterGroup) {
	app := api.App.UserApi
	r.POST("user", app.UserCreateView)

}

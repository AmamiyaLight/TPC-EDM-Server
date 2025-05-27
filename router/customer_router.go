package router

import (
	"TPC-EDM-Server/api"
	"github.com/gin-gonic/gin"
)

func CustomerRouter(r *gin.RouterGroup) {
	app := api.App.CustomerApi
	r.POST("Customer", app.CustomerInsertView)
	r.GET("Customer", app.CustomerListView)
}

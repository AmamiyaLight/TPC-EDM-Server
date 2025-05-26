package router

import (
	"TPC-H-EDM-Server/api"
	"github.com/gin-gonic/gin"
)

func OrdersRouter(r *gin.RouterGroup) {
	app := api.App.OrdersApi
	r.POST("order", app.OrderInsertView)
	r.GET("order", app.OrdersListView)
}

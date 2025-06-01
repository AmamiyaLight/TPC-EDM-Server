package router

import (
	"TPC-EDM-Server/api"
	"github.com/gin-gonic/gin"
)

func TpccRouter(r *gin.RouterGroup) {
	app := api.App.TpccApi
	r.POST("tpcc1", app.NewOrderView)
	r.POST("tpcc2", app.PaymentView)
}

package router

import (
	"TPC-EDM-Server/api"
	"github.com/gin-gonic/gin"
)

func SupplierRouter(r *gin.RouterGroup) {
	app := api.App.SupplierApi
	r.POST("Supplier", app.SupplierInsertView)
	r.GET("Supplier", app.SupplierListView)
}

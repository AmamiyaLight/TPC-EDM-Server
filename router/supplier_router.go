package router

import (
	"TPC-EDM-Server/api"
	"github.com/gin-gonic/gin"
)

func SupplierRouter(r *gin.RouterGroup) {
	app := api.App.SupplierApi
	r.POST("supplier", app.SupplierInsertView)
	r.GET("supplier", app.SupplierListView)
	r.GET("supplier/download", app.SupplierDownloadView)
	r.GET("supplier/top", app.TopSupplierView)
}

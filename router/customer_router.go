package router

import (
	"TPC-EDM-Server/api"
	"github.com/gin-gonic/gin"
)

func CustomerRouter(r *gin.RouterGroup) {
	app := api.App.CustomerApi
	r.POST("customer", app.CustomerInsertView)
	r.GET("customer", app.CustomerListView)
	r.GET("customer/download", app.CustomerDownloadView)
	r.GET("customer/search", app.CustomerSearchView)
}

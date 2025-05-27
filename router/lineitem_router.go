package router

import (
	"TPC-EDM-Server/api"
	"github.com/gin-gonic/gin"
)

func LineItemRouter(r *gin.RouterGroup) {
	app := api.App.LineItemApi
	r.POST("lineitem", app.LineItemInsertView)
	r.GET("lineitem", app.LineItemListView)
}

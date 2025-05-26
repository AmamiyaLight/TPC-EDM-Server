package router

import (
	"TPC-H-EDM-Server/api"
	"github.com/gin-gonic/gin"
)

func LineItemRouter(r *gin.RouterGroup) {
	app := api.App.LineItemApi
	r.POST("lineitem", app.LineItemInsertView)
	r.GET("lineitem", app.LineItemListView)
}

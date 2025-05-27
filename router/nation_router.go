package router

import (
	"TPC-EDM-Server/api"
	"github.com/gin-gonic/gin"
)

func NationRouter(r *gin.RouterGroup) {
	app := api.App.NationApi
	r.POST("Nation", app.NationInsertView)
	r.GET("Nation", app.NationListView)
}

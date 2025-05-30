package router

import (
	"TPC-EDM-Server/api"
	"github.com/gin-gonic/gin"
)

func NationRouter(r *gin.RouterGroup) {
	app := api.App.NationApi
	r.POST("nation", app.NationInsertView)
	r.GET("nation", app.NationListView)
	r.GET("nation/download", app.NationDownloadView)
}

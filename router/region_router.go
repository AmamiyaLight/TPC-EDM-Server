package router

import (
	"TPC-EDM-Server/api"
	"github.com/gin-gonic/gin"
)

func RegionRouter(r *gin.RouterGroup) {
	app := api.App.RegionApi
	r.POST("region", app.RegionInsertView)
	r.GET("region", app.RegionListView)
	r.GET("region/download", app.RegionDownloadView)
}

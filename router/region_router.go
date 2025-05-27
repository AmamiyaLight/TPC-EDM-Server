package router

import (
	"TPC-EDM-Server/api"
	"github.com/gin-gonic/gin"
)

func RegionRouter(r *gin.RouterGroup) {
	app := api.App.RegionApi
	r.POST("Region", app.RegionInsertView)
	r.GET("Region", app.RegionListView)
}

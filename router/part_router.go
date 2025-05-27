package router

import (
	"TPC-EDM-Server/api"
	"github.com/gin-gonic/gin"
)

func PartRouter(r *gin.RouterGroup) {
	app := api.App.PartApi
	r.POST("Part", app.PartInsertView)
	r.GET("Part", app.PartListView)
}

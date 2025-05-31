package router

import (
	"TPC-EDM-Server/api"
	"github.com/gin-gonic/gin"
)

func PartRouter(r *gin.RouterGroup) {
	app := api.App.PartApi
	r.POST("part", app.PartInsertView)
	r.GET("part", app.PartListView)
	r.GET("part/download", app.PartDownloadView)
	r.GET("part/promo", app.PartPromoView)
}

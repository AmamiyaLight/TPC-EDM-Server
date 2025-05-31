package router

import (
	"TPC-EDM-Server/api"
	"github.com/gin-gonic/gin"
)

func PartSuppRouter(r *gin.RouterGroup) {
	app := api.App.PartSuppApi
	r.POST("part_supp", app.PartSuppInsertView)
	r.GET("part_supp", app.PartSuppListView)
	r.GET("part_supp/download", app.PartSuppDownloadView)
	r.GET("part_supp/relation", app.PartSuppRelationView)
}

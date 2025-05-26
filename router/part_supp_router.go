package router

import (
	"TPC-H-EDM-Server/api"
	"github.com/gin-gonic/gin"
)

func PartSuppRouter(r *gin.RouterGroup) {
	app := api.App.PartSuppApi
	r.POST("part_supp", app.PartSuppInsertView)
	r.GET("part_supp", app.PartSuppListView)
}

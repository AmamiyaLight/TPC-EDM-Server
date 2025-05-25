package router

import (
	"TPC-H-EDM-Server/global"
	"github.com/gin-gonic/gin"
)

func Run() {
	gin.SetMode(global.Config.System.GinMode)
	r := gin.Default()

	r.Static("/uploads", "uploads")

	nr := r.Group("/api")
	addr := global.Config.System.Addr()
	UserRouter(nr)
	r.Run(addr)
}

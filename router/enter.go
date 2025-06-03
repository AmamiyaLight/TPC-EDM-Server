package router

import (
	"TPC-EDM-Server/global"
	"TPC-EDM-Server/middleware"
	"github.com/gin-gonic/gin"
)

func Run() {
	gin.SetMode(global.Config.System.GinMode)
	r := gin.Default()
	r.Use(middleware.Cors())
	r.Use(middleware.AuthMiddleware)
	nr := r.Group("/api")
	addr := global.Config.System.Addr()
	UserRouter(nr)
	OrdersRouter(nr)
	PartSuppRouter(nr)
	LineItemRouter(nr)
	NationRouter(nr)
	CustomerRouter(nr)
	PartRouter(nr)
	RegionRouter(nr)
	SupplierRouter(nr)
	TpccRouter(nr)
	SysRouter(nr)
	r.Run(addr)
}

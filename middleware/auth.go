package middleware

import (
	"TPC-EDM-Server/common/res"
	"TPC-EDM-Server/models/enum"
	"TPC-EDM-Server/utils/jwts"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	cla, err := jwts.ParseTokenByGin(c)
	if err != nil {
		res.FailWithError(err, c)
		c.Abort()
		return
	}
	c.Set("claims", cla)
	return
}

func AdminMiddleware(c *gin.Context) {
	cla, err := jwts.ParseTokenByGin(c)
	if err != nil {
		res.FailWithError(err, c)
		c.Abort()
		return
	}
	if cla.Role != enum.AdminRole {
		res.FailWithMsg("权限错误", c)
		c.Abort()
		return
	}

	c.Set("claims", cla)
	return
}

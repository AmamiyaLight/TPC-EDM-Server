package user_api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserApi struct{}

func (UserApi) UserCreateView(c *gin.Context) {
	logrus.Infof("测试通过")
	return
}

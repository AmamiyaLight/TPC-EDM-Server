package user_api

import (
	"TPC-EDM-Server/common/res"
	"TPC-EDM-Server/global"
	"TPC-EDM-Server/models"
	"TPC-EDM-Server/models/enum"
	"github.com/gin-gonic/gin"
)

func (UserApi) UserDeleteView(c *gin.Context) {
	var cr models.IDRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	var user models.UserModel
	err = global.DB.Take(&user, cr.ID).Error
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}
	if user.Role == enum.AdminRole {
		res.FailWithMsg("不可以删除管理员用户", c)
		return
	}

	err = global.DB.Delete(&user).Error
	if err != nil {
		res.FailWithMsg("用户删除失败", c)
		return
	}
	res.OkWithMsg("用户删除成功", c)
}

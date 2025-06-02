package user_api

import (
	"TPC-EDM-Server/common/res"
	"TPC-EDM-Server/global"
	"TPC-EDM-Server/models"
	"TPC-EDM-Server/utils/jwts"
	"TPC-EDM-Server/utils/pwd"
	"github.com/gin-gonic/gin"
)

type PwdLoginRequest struct {
	Val string `json:"val" binding:"required"`
	Pwd string `json:"pwd" binding:"required"`
}

func (UserApi) PwdLoginView(c *gin.Context) {
	var cr PwdLoginRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	if cr.Val == "" {
		res.FailWithMsg("用户名不能为空", c)
		return
	}
	if cr.Pwd == "" {
		res.FailWithMsg("密码不能为空", c)
		return
	}
	var user models.UserModel

	err = global.DB.Where("username = ?", cr.Val).Take(&user).Error
	if err != nil {
		res.FailWithMsg("用户名不存在", c)
		return
	}

	if !pwd.CompareHashAndPassword(user.Password, cr.Pwd) {
		res.FailWithMsg("密码错误", c)
		return
	}
	token, err := jwts.GenerateJWT(jwts.Claims{
		user.ID,
		user.Username,
		user.Role,
	})
	if err != nil {
		res.FailWithMsg("登录凭证生成失败", c)
		return
	}

	res.OkWithData(token, c)

}

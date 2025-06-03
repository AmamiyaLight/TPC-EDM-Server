package user_api

import (
	"TPC-EDM-Server/common/res"
	"TPC-EDM-Server/global"
	"TPC-EDM-Server/models"
	"TPC-EDM-Server/models/enum"
	"TPC-EDM-Server/utils/jwts"
	"github.com/gin-gonic/gin"
)

type UserDetailResponse struct {
	ID       uint          `json:"id"`
	Username string        `json:"username"`
	Role     enum.RoleType `json:"role"`
}

func (UserApi) UserDetailView(c *gin.Context) {
	claims := jwts.GetClaims(c)
	var user models.UserModel
	err := global.DB.Take(&user, claims.UserID).Error
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}
	var data = UserDetailResponse{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}
	res.OkWithData(data, c)
}

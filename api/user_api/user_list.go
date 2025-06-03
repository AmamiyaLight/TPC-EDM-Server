package user_api

import (
	"TPC-EDM-Server/common"
	"TPC-EDM-Server/common/res"
	"TPC-EDM-Server/models"
	"TPC-EDM-Server/models/enum"
	"github.com/gin-gonic/gin"
	"time"
)

type UserListRequest struct {
	common.PageInfo
}

type UserListResponse struct {
	UserId    uint          `json:"user_id"`
	Username  string        `json:"username"`
	Role      enum.RoleType `json:"role"`
	CreatedAt time.Time     `json:"created_at"`
}

func (UserApi) UserListView(c *gin.Context) {
	var cr UserListRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	_list, count, _ := common.ListQuery(models.UserModel{}, common.Options{
		PageInfo: cr.PageInfo,
	})

	list := make([]UserListResponse, 0)
	for _, v := range _list {
		item := UserListResponse{
			UserId:    v.ID,
			Username:  v.Username,
			Role:      v.Role,
			CreatedAt: v.CreatedAt,
		}
		list = append(list, item)
	}
	res.OkWithList(list, count, c)
}

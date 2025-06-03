package system_api

import (
	"TPC-EDM-Server/common/res"
	"TPC-EDM-Server/conf"
	"TPC-EDM-Server/core"
	"TPC-EDM-Server/global"
	"github.com/gin-gonic/gin"
)

type SystemApi struct {
}

type SystemInfoResponse struct {
	User        string `json:"user"`
	Password    string `json:"password"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	DB          string `json:"DB"`
	Source      string `json:"source"`
	MaxConn     int    `json:"maxConn"`
	MaxIdle     int    `json:"maxIdle"`
	MaxLifeTime int    `json:"maxLifeTime"`
}

func (SystemApi) SystemInfoView(c *gin.Context) {
	res.OkWithData(SystemInfoResponse{
		User:        global.Config.DB[0].User,
		Password:    "******",
		Host:        global.Config.DB[0].Host,
		Port:        global.Config.DB[0].Port,
		DB:          global.Config.DB[0].DB,
		Source:      global.Config.DB[0].Source,
		MaxConn:     global.Config.DB[0].MaxConn,
		MaxIdle:     global.Config.DB[0].MaxIdle,
		MaxLifeTime: global.Config.DB[0].MaxLifeTime,
	}, c)
}

type SystemInfoRequest struct {
	Name string `uri:"name"`
}

func (SystemApi) SystemUpdateView(c *gin.Context) {

	var cr SystemInfoRequest
	err := c.ShouldBindUri(&cr)

	if err != nil {
		res.FailWithError(err, c)
		return
	}
	var rep any

	switch cr.Name {
	case "db":
		var data conf.DB
		err = c.ShouldBind(&data)
		rep = data
	default:
		res.FailWithMsg("不存在的配置", c)
		return
	}
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	switch s := rep.(type) {
	case conf.DB:
		if s.Password == "******" {
			s.Password = global.Config.DB[0].Password
		}
		global.Config.DB[0] = s
		return
	}

	core.SetConf()

	res.OkWithMsg("更新成功", c)
	return
}

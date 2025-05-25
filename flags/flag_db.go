package flags

import (
	"TPC-H-EDM-Server/global"
	"github.com/sirupsen/logrus"
)

func FlagDB() {
	err := global.DB.AutoMigrate()
	if err != nil {
		logrus.Errorf("数据库迁移失败 %s\n", err)
		return
	}
	logrus.Infoln("数据库迁移成功")
}

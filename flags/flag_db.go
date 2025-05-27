package flags

import (
	"TPC-EDM-Server/global"
	"TPC-EDM-Server/models"
	"github.com/sirupsen/logrus"
)

func FlagDB() {
	err := global.DB.AutoMigrate(
		&models.UserModel{},
		&models.CustomerModel{},
		&models.LineItemModel{},
		&models.NationModel{},
		&models.OrdersModel{},
		&models.PartModel{},
		&models.PartSuppModel{},
		&models.RegionModel{},
		&models.SupplierModel{},
	)
	if err != nil {
		logrus.Errorf("数据库迁移失败 %s\n", err)
		return
	}
	logrus.Infoln("数据库迁移成功")
}

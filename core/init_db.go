package core

import (
	"TPC-EDM-Server/global"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"time"
)

func InitDB() *gorm.DB {
	if len(global.Config.DB) == 0 {
		logrus.Fatalf("未配置数据库")
	}

	dc := global.Config.DB[0]

	// TODO:PGSQL支持

	db, err := gorm.Open(mysql.Open(dc.DSN()), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		logrus.Fatalf("数据库连接失败 %s\n", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		logrus.Fatal(err)
	}
	sqlDB.SetMaxIdleConns(global.Config.DB[0].MaxIdle)
	sqlDB.SetMaxOpenConns(global.Config.DB[0].MaxConn)
	sqlDB.SetConnMaxLifetime(time.Hour * time.Duration(global.Config.DB[0].MaxLifeTime))
	logrus.Infoln("数据库连接成功")
	if len(global.Config.DB) > 1 {
		var readList []gorm.Dialector
		for _, d := range global.Config.DB[1:] {
			readList = append(readList, mysql.Open(d.DSN()))
		}
		// 注册读写分离
		err = db.Use(dbresolver.Register(dbresolver.Config{
			Sources:           []gorm.Dialector{mysql.Open(dc.DSN())}, //写
			Replicas:          readList,                               //读
			Policy:            dbresolver.RandomPolicy{},
			TraceResolverMode: true,
		}))
		if err != nil {
			logrus.Fatalf("读写配置出错 %s\n", err)
		}
	}

	return db
}

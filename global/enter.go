package global

import (
	"TPC-H-EDM-Server/conf"
	"gorm.io/gorm"
)

var (
	Config *conf.Config
	DB     *gorm.DB
)

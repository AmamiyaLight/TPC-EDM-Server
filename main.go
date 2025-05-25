package main

import (
	"TPC-H-EDM-Server/core"
	"TPC-H-EDM-Server/flags"
	"TPC-H-EDM-Server/global"
	"TPC-H-EDM-Server/router"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.Initlogrus()
	global.DB = core.InitDB()
	flags.Run()
	router.Run()
}

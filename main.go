package main

import (
	"TPC-EDM-Server/core"
	"TPC-EDM-Server/flags"
	"TPC-EDM-Server/global"
	"TPC-EDM-Server/router"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	core.Initlogrus()
	global.DB = core.InitDB()
	flags.Run()
	router.Run()
}

package flags

import (
	"TPC-EDM-Server/service"
	"flag"
	"os"
)

type Options struct {
	File    string
	DB      bool
	Version bool
	Type    string
	Sub     string
	TPCH    bool
}

var FlagOptions = new(Options)

func Parse() {
	flag.StringVar(&FlagOptions.File, "f", "settings.yaml", "configuration file")
	flag.BoolVar(&FlagOptions.DB, "db", false, "use db")
	flag.BoolVar(&FlagOptions.Version, "v", false, "show version")
	flag.StringVar(&FlagOptions.Type, "t", "", "type")
	flag.StringVar(&FlagOptions.Sub, "s", "", "sub")
	flag.BoolVar(&FlagOptions.TPCH, "tpch", false, "test tpc-h")
	flag.Parse()
}

func Run() {
	if FlagOptions.DB {
		//迁移数据库
		FlagDB()
		os.Exit(0)
	}
	if FlagOptions.TPCH {
		service.TpchTest()
		os.Exit(0)
	}
}

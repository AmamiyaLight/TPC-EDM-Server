package flags

import (
	"flag"
	"os"
)

type Options struct {
	File    string
	DB      bool
	Version bool
	Type    string
	Sub     string
	ES      bool
}

var FlagOptions = new(Options)

func Parse() {
	flag.StringVar(&FlagOptions.File, "f", "settings.yaml", "configuration file")
	flag.BoolVar(&FlagOptions.DB, "db", false, "use db")
	flag.BoolVar(&FlagOptions.Version, "v", false, "show version")
	flag.StringVar(&FlagOptions.Type, "t", "", "type")
	flag.StringVar(&FlagOptions.Sub, "s", "", "sub")
	flag.Parse()
}

func Run() {
	if FlagOptions.DB {
		//迁移数据库
		FlagDB()
		os.Exit(0)
	}
}

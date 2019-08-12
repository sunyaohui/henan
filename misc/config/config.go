package config

import (
	"flag"
	"runtime"

	"github.com/larspensjo/config"
	"github.com/wonderivan/logger"
)

var (
	configFile = flag.String("configfile", "config.ini", "General configuration file")
)

//topic list
var CONFIG = make(map[string]string)

func InitConfig() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	//set config file std
	cfg, err := config.ReadDefault(*configFile)
	if err != nil {
		logger.Warn("Fail to find", *configFile, err)
	}
	//set config file std End

	//Initialized topic from the configuration
	if cfg.HasSection("config") {
		section, err := cfg.SectionOptions("config")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("config", v)
				if err == nil {
					CONFIG[v] = options
				}
			}
		}
	}
	//Initialized topic from the configuration END
}

package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

func LoadConfig() {
	if Conf == nil {
		Conf = &Config{}
	}
	mt, err := toml.DecodeFile(RootPath+"/config/config.toml", Conf)
	if err != nil {
		panic(err)
	}
	fmt.Println(mt)
	if Conf.GoSwitchPath == "" && RootPath != "" {
		Conf.GoSwitchPath = RootPath
	}
}

func InitConfigFile() {
	_, err := os.Create(RootPath + "/config/config.toml")
	if err != nil {
		panic(err)
	}
}

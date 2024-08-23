package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

func LoadConfig() {
	if Conf == nil {
		Conf = &Config{}
	}
	configFilePath := filepath.Join(RootPath, "config")
	_, err := toml.DecodeFile(fmt.Sprintf("%s%s%s", configFilePath, string(os.PathSeparator), "config.toml"), Conf)
	if err != nil {
		panic(err)
	}
	if Conf.GoSwitchPath == "" && RootPath != "" {
		Conf.GoSwitchPath = RootPath
		Conf.SaveConfig()
	}
}

func InitConfigFile() {

	if exists, create := ExistsPath(RootPath); !exists && !create {
		panic("RootPath not exists")
	}
	if err := GlobalSetPermissions.SetHiddenAttribute(RootPath); err != nil {
		panic("RootPath SetHiddenAttribute failed " + err.Error())
	}

	if exists, create := ExistsPath(GoEnvFilePath); !exists && !create {
		panic("GoEnvFilePath not exists")
	}

	configPath := filepath.Join(RootPath, "config")
	if exists, create := ExistsPath(configPath); !exists && !create {
		panic("configPath not exists")
	}

	if exists, create := FileExists(fmt.Sprintf("%s%s%s", configPath, string(os.PathSeparator), "config.toml")); !exists && !create {
		panic("config file not exists")
	}

	if exists, create := FileExists(fmt.Sprintf("%s%s%s", GoEnvFilePath, string(os.PathSeparator), "system")); !exists && !create {
		panic("system env file not exists")
	}

}

func InitSystemVars() {

	systemOs := runtime.GOOS
	switch systemOs {
	case "linux":
		SystemEnv = Linux
		RootPath = filepath.Join(LinuxGoPath, GoSwitchDir)

	case "windows":
		SystemEnv = Windows
		GoSwitchDir = "go-switch"
		RootPath = filepath.Join(WindowsGoPath, GoSwitchDir)
	case "darwin":
		SystemEnv = Mac
		RootPath = filepath.Join(MacGoPath + GoSwitchDir)
	}
	GosPath = filepath.Join(RootPath, SaveGoDir)
	TempUnzipPath = filepath.Join(GosPath, UnzipGoDir)
	SystemArch = runtime.GOARCH
	GoEnvFilePath = filepath.Join(RootPath, "environment")
}

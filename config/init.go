package config

import (
	"fmt"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/xulimeng/go-switch/utils"
)

func LoadConfig() {
	if Conf == nil {
		Conf = &Config{}
	}
	configFilePath := filepath.Join(RootPath+"config", "config.toml")
	fmt.Println("---------------configFilePath: ", configFilePath)
	_, err := toml.DecodeFile(configFilePath, Conf)
	if err != nil {
		panic(err)
	}
	if Conf.GoSwitchPath == "" && RootPath != "" {
		Conf.GoSwitchPath = RootPath
		Conf.SaveConfig()
	}
}

func InitConfigFile() {

	if exists, create := utils.FileExists(filepath.Join(RootPath, "config", "config.toml")); !exists && !create {
		panic("config file not exists")
	}

	if exists, create := utils.FileExists(filepath.Join(GoEnvFilePath, "system")); !exists && !create {
		panic("system env file not exists")
	}

}

func InitSystemVars() {

	os := runtime.GOOS
	switch os {
	case "linux":
		SystemEnv = Linux
		RootPath = LinuxGoPath + GoSwitchDir
		GosPath = RootPath + "/" + SaveGoDir
		TempUnzipPath = GosPath + "/" + UnzipGoDir

	case "windows":
		SystemEnv = Windows
		userNameCurr, err := user.Current()
		if err != nil {
			panic(err)
		}
		RootPath = WindowsGoPath + userNameCurr.Username + "\\" + GoSwitchDir
		GosPath = RootPath + "\\" + SaveGoDir
		TempUnzipPath = GosPath + "\\" + UnzipGoDir
	case "darwin":
		SystemEnv = Mac
		RootPath = MacGoPath + GoSwitchDir
		GosPath = RootPath + "/" + SaveGoDir
		TempUnzipPath = GosPath + "/" + UnzipGoDir
	}
	SystemArch = runtime.GOARCH
	GoEnvFilePath = filepath.Join(RootPath, "environment")
}

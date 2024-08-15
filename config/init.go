package config

import (
	"fmt"
	"os"
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
	configFilePath := filepath.Join(RootPath, "config")
	fmt.Println("---------------configFilePath: ", configFilePath)
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

	if exists, create := utils.FileExists(fmt.Sprintf("%s%s%s", filepath.Join(RootPath, "config"), string(os.PathSeparator), "config.toml")); !exists && !create {
		panic("config file not exists")
	}

	if exists, create := utils.FileExists(fmt.Sprintf("%s%s%s", GoEnvFilePath, string(os.PathSeparator), "system")); !exists && !create {
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

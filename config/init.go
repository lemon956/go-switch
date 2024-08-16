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
	fmt.Println("configPath", fmt.Sprintf("%s%s%s", filepath.Join(RootPath, "config"), string(os.PathSeparator), "config.toml"))
	if exists, create := FileExists(fmt.Sprintf("%s%s%s", filepath.Join(RootPath, "config"), string(os.PathSeparator), "config.toml")); !exists && !create {
		panic("config file not exists")
	}
	fmt.Println("GoEnvFilePath", fmt.Sprintf("%s%s%s", GoEnvFilePath, string(os.PathSeparator), "system"))
	if exists, create := FileExists(fmt.Sprintf("%s%s%s", GoEnvFilePath, string(os.PathSeparator), "system")); !exists && !create {
		panic("system env file not exists")
	}

}

func InitSystemVars() {

	os := runtime.GOOS
	switch os {
	case "linux":
		SystemEnv = Linux
		RootPath = LinuxGoPath + GoSwitchDir

	case "windows":
		SystemEnv = Windows
		RootPath = WindowsGoPath + "\\" + GoSwitchDir
	case "darwin":
		SystemEnv = Mac
		RootPath = MacGoPath + GoSwitchDir
	}
	GosPath = filepath.Join(RootPath, SaveGoDir)
	TempUnzipPath = filepath.Join(GosPath, UnzipGoDir)
	SystemArch = runtime.GOARCH
	GoEnvFilePath = filepath.Join(RootPath, "environment")
}

package config

import (
	"fmt"
	"os/user"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/xulimeng/go-switch/utils"
)

func LoadConfig() {
	if Conf == nil {
		Conf = &Config{}
	}
	configFilePath := RootPath + "/config/config.toml"
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
	if exists, create := utils.FileExists(RootPath + "/config/config.toml"); !exists && !create {
		panic("config file not exists")
	}
	if exists, create := utils.FileExists(GoEnvFilePath + "/system"); !exists && !create {
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
	GoEnvFilePath = ConnectPathWithEnv(SystemEnv, RootPath, []string{"environment"})
}

// ConnectPathWithEnv 根据不同系统环境拼接路径
func ConnectPathWithEnv(env Env, basePath string, connectPaths []string) string {
	if env == Linux || env == Mac {
		return fmt.Sprintf("%s/%s", basePath, strings.Join(connectPaths, "/"))
	} else if env == Windows {
		return fmt.Sprintf("%s\\%s", basePath, strings.Join(connectPaths, "\\"))
	} else {
		return ""
	}
}

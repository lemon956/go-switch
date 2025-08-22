package features

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/xulimeng/go-switch/config"
	"github.com/xulimeng/go-switch/helper"
)

func LoadConfig() {
	if config.Conf == nil {
		config.Conf = &config.Config{}
	}
	configFilePath := filepath.Join(config.RootPath, "config")
	_, err := toml.DecodeFile(fmt.Sprintf("%s%s%s", configFilePath, string(os.PathSeparator), "config.toml"), config.Conf)
	if err != nil {
		panic(err)
	}
	if config.Conf.GoSwitchPath == "" && config.RootPath != "" {
		config.Conf.GoSwitchPath = config.RootPath
		config.Conf.SaveConfig()
	}
	// 初始化GOPATH，只在第一次初始化时设置
	if config.Conf.GoPath == "" && config.GoPathDirPath != "" {
		config.Conf.GoPath = config.GoPathDirPath
		// 一次性将GOPATH写入system文件
		initGoPathToSystem()
		config.Conf.SaveConfig()
	}
}

func InitConfigFile() {

	if exists, create := helper.ExistsPath(config.RootPath); !exists && !create {
		panic("RootPath not exists")
	}
	if err := helper.GlobalSetPermissions.SetHiddenAttribute(config.RootPath); err != nil {
		panic("RootPath SetHiddenAttribute failed " + err.Error())
	}

	if exists, create := helper.ExistsPath(config.GoEnvFilePath); !exists && !create {
		panic("GoEnvFilePath not exists")
	}

	// 创建GOPATH目录
	if exists, create := helper.ExistsPath(config.GoPathDirPath); !exists && !create {
		panic("GoPathDirPath not exists")
	}

	configPath := filepath.Join(config.RootPath, "config")
	if exists, create := helper.ExistsPath(configPath); !exists && !create {
		panic("configPath not exists")
	}

	if exists, create := helper.FileExists(fmt.Sprintf("%s%s%s", configPath, string(os.PathSeparator), "config.toml")); !exists && !create {
		panic("config file not exists")
	}

	if exists, create := helper.FileExists(fmt.Sprintf("%s%s%s", config.GoEnvFilePath, string(os.PathSeparator), "system")); !exists && !create {
		panic("system env file not exists")
	}

}

func InitSystemVars() {

	systemOs := runtime.GOOS
	switch systemOs {
	case "linux":
		config.SystemEnv = config.Linux
		config.RootPath = filepath.Join(config.LinuxGoPath, config.GoSwitchDir)

	case "windows":
		config.SystemEnv = config.Windows
		config.GoSwitchDir = "go-switch"
		config.RootPath = filepath.Join(config.WindowsGoPath, config.GoSwitchDir)
	case "darwin":
		config.SystemEnv = config.Mac
		config.RootPath = filepath.Join(config.MacGoPath + config.GoSwitchDir)
	}
	config.GosPath = filepath.Join(config.RootPath, config.SaveGoDir)
	config.TempUnzipPath = filepath.Join(config.GosPath, config.UnzipGoDir)
	config.SystemArch = runtime.GOARCH
	config.GoEnvFilePath = filepath.Join(config.RootPath, "environment")
	config.GoPathDirPath = filepath.Join(config.RootPath, config.GoPathDir)
}

// initGoPathToSystem 在初始化时一次性将GOPATH写入system文件
func initGoPathToSystem() {
	if config.Conf.GoPath == "" {
		return
	}

	goEnvFilePath := filepath.Join(config.GoEnvFilePath, "system")
	goPathCmd := fmt.Sprintf("export GOPATH=%s", config.Conf.GoPath)

	// 检查system文件是否存在，如果不存在则创建
	if _, err := os.Stat(goEnvFilePath); os.IsNotExist(err) {
		file, err := os.Create(goEnvFilePath)
		if err != nil {
			fmt.Printf("Failed to create system file: %v\n", err)
			return
		}
		file.Close()
	}

	// 检查GOPATH是否已经存在
	file, err := os.Open(goEnvFilePath)
	if err != nil {
		fmt.Printf("Failed to open system file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	found := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "export GOPATH=") {
			found = true
			break
		}
	}

	// 如果没有找到GOPATH配置，则添加
	if !found {
		file, err := os.OpenFile(goEnvFilePath, os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Printf("Failed to open system file for writing: %v\n", err)
			return
		}
		defer file.Close()

		if _, err := file.WriteString(goPathCmd + "\n"); err != nil {
			fmt.Printf("Failed to write GOPATH to system file: %v\n", err)
		} else {
			fmt.Printf("Added GOPATH configuration to system file: %s\n", goPathCmd)
		}
	}
}

// initGoSwitch 初始化 go-switch 环境
func InitGoSwitch() {
	InitSystemVars()
	InitConfigFile()
	LoadConfig()
	fmt.Println("go-switch 环境初始化完成！")
}

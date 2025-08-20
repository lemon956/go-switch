package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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
	// 初始化GOPATH，只在第一次初始化时设置
	if Conf.GoPath == "" && GoPathDirPath != "" {
		Conf.GoPath = GoPathDirPath
		// 一次性将GOPATH写入system文件
		initGoPathToSystem()
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

	// 创建GOPATH目录
	if exists, create := ExistsPath(GoPathDirPath); !exists && !create {
		panic("GoPathDirPath not exists")
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
	GoPathDirPath = filepath.Join(RootPath, GoPathDir)
}

// initGoPathToSystem 在初始化时一次性将GOPATH写入system文件
func initGoPathToSystem() {
	if Conf.GoPath == "" {
		return
	}

	goEnvFilePath := filepath.Join(GoEnvFilePath, "system")
	goPathCmd := fmt.Sprintf("export GOPATH=%s", Conf.GoPath)

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

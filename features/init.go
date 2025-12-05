package features

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/xulimeng/go-switch/config"
	"github.com/xulimeng/go-switch/helper"
)

// LoadConfig 加载配置文件到全局 config.Conf
// 若配置文件不存在或内容不合法，返回 error 由调用方统一处理
func LoadConfig() error {
	if config.Conf == nil {
		config.Conf = &config.Config{}
	}
	configDir := filepath.Join(config.RootPath, "config")
	configFile := filepath.Join(configDir, "config.toml")

	if _, err := os.Stat(configFile); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Config file does not exist: %s, please run 'goswitch init' first", configFile)
		}
		return fmt.Errorf("Failed to check config file: %w", err)
	}

	if _, err := toml.DecodeFile(configFile, config.Conf); err != nil {
		return fmt.Errorf("Failed to load config file %s: %w", configFile, err)
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
	return nil
}

// InitConfigFile 初始化 go-switch 相关目录与基础文件，仅在 init 命令中调用
func InitConfigFile() error {
	if exists, create := helper.ExistsPath(config.RootPath); !exists && !create {
		return fmt.Errorf("RootPath not exists: %s", config.RootPath)
	}
	if err := helper.GlobalSetPermissions.SetHiddenAttribute(config.RootPath); err != nil {
		return fmt.Errorf("RootPath SetHiddenAttribute failed: %w", err)
	}

	if exists, create := helper.ExistsPath(config.GoEnvFilePath); !exists && !create {
		return fmt.Errorf("GoEnvFilePath not exists: %s", config.GoEnvFilePath)
	}

	// 创建GOPATH目录
	if exists, create := helper.ExistsPath(config.GoPathDirPath); !exists && !create {
		return fmt.Errorf("GoPathDirPath not exists: %s", config.GoPathDirPath)
	}

	configPath := filepath.Join(config.RootPath, "config")
	if exists, create := helper.ExistsPath(configPath); !exists && !create {
		return fmt.Errorf("configPath not exists: %s", configPath)
	}

	if exists, create := helper.FileExists(filepath.Join(configPath, "config.toml")); !exists && !create {
		return errors.New("config file not exists")
	}

	if exists, create := helper.FileExists(filepath.Join(config.GoEnvFilePath, "system")); !exists && !create {
		return errors.New("system env file not exists")
	}

	// Create fish environment file as well
	if exists, create := helper.FileExists(filepath.Join(config.GoEnvFilePath, "system.fish")); !exists && !create {
		return errors.New("system.fish env file not exists")
	}
	return nil
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

// initGoPathToSystem writes GOPATH to the system environment files during initialization
// It writes to both bash/zsh format (system) and fish format (system.fish)
func initGoPathToSystem() {
	if config.Conf.GoPath == "" {
		return
	}

	// Write GOPATH for bash/zsh (export syntax)
	writeGoPathForShell("bash", config.Conf.GoPath)

	// Write GOPATH for fish (set -gx syntax)
	writeGoPathForShell("fish", config.Conf.GoPath)
}

// writeGoPathForShell writes GOPATH to the environment file for the specified shell
func writeGoPathForShell(shell, goPath string) {
	var goEnvFilePath string
	var goPathCmd string
	var checkPrefix string

	if shell == "fish" {
		goEnvFilePath = filepath.Join(config.GoEnvFilePath, "system.fish")
		goPathCmd = fmt.Sprintf("set -gx GOPATH %s", goPath)
		checkPrefix = "set -gx GOPATH"
	} else {
		goEnvFilePath = filepath.Join(config.GoEnvFilePath, "system")
		goPathCmd = fmt.Sprintf("export GOPATH=%s", goPath)
		checkPrefix = "export GOPATH="
	}

	// Check if system file exists, create if not
	if _, err := os.Stat(goEnvFilePath); os.IsNotExist(err) {
		file, err := os.Create(goEnvFilePath)
		if err != nil {
			fmt.Printf("Failed to create %s file: %v\n", shell, err)
			return
		}
		file.Close()
	}

	// Check if GOPATH already exists
	file, err := os.Open(goEnvFilePath)
	if err != nil {
		fmt.Printf("Failed to open %s file: %v\n", shell, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	found := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, checkPrefix) {
			found = true
			break
		}
	}

	// If GOPATH config not found, add it
	if !found {
		file, err := os.OpenFile(goEnvFilePath, os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Printf("Failed to open %s file for writing: %v\n", shell, err)
			return
		}
		defer file.Close()

		if _, err := file.WriteString(goPathCmd + "\n"); err != nil {
			fmt.Printf("Failed to write GOPATH to %s file: %v\n", shell, err)
		} else {
			fmt.Printf("Added GOPATH configuration to %s file: %s\n", shell, goPathCmd)
		}
	}

	// Ensure PATH includes GOPATH/bin for this shell (during initialization)
	var pathPrefix, pathCmd, goPathSegment string
	if shell == "fish" {
		pathPrefix = "set -gx PATH"
		goPathSegment = "$GOPATH/bin"
		pathCmd = "set -gx PATH $GOPATH/bin $PATH"
	} else {
		pathPrefix = "export PATH="
		goPathSegment = "$GOPATH/bin"
		pathCmd = "export PATH=$GOPATH/bin:$PATH"
	}

	// Reopen file to check for existing PATH with GOPATH/bin
	file2, err := os.Open(goEnvFilePath)
	if err != nil {
		fmt.Printf("Failed to open %s file for PATH check: %v\n", shell, err)
		return
	}
	defer file2.Close()

	scanner2 := bufio.NewScanner(file2)
	pathWithGoPath := false
	for scanner2.Scan() {
		line := strings.TrimSpace(scanner2.Text())
		if strings.HasPrefix(line, pathPrefix) && strings.Contains(line, goPathSegment) {
			pathWithGoPath = true
			break
		}
	}

	if !pathWithGoPath {
		file2, err := os.OpenFile(goEnvFilePath, os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Printf("Failed to open %s file for PATH writing: %v\n", shell, err)
			return
		}
		defer file2.Close()

		if _, err := file2.WriteString(pathCmd + "\n"); err != nil {
			fmt.Printf("Failed to write PATH with GOPATH to %s file: %v\n", shell, err)
		} else {
			fmt.Printf("Added PATH with GOPATH/bin to %s file: %s\n", shell, pathCmd)
		}
	}
}

// InitGoSwitch 初始化 go-switch 环境
func InitGoSwitch() error {
	InitSystemVars()
	if err := InitConfigFile(); err != nil {
		return fmt.Errorf("Failed to initialize config: %w", err)
	}
	if err := LoadConfig(); err != nil {
		return fmt.Errorf("Failed to load config: %w", err)
	}
	fmt.Println("go-switch environment initialization completed!")
	return nil
}

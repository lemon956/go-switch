package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xulimeng/go-switch/config"
	"github.com/xulimeng/go-switch/features"
	"github.com/xulimeng/go-switch/helper"
)

func PrintHelp() {
	fmt.Println("Usage: goswitch -cmd <command>")

	fmt.Println("\ngoswitch is the Go Version Manager")

	fmt.Println(`
Command:
	help	- Show this help message
	init	- Initialize goswitch environment
	install	- Install to go version
	switch	- Choose go version
	list	- List all installed go versions
	listall - List all available go versions
	delete	- Delete go version
	env 	- Show goswitch environment
	`)
}

// ensureInitialized 确保系统已初始化（除了init命令外的其他命令需要）
func ensureInitialized() {
	features.InitSystemVars()
	// 通过配置文件是否存在判断是否已经初始化，而不是仅依赖 RootPath 是否为空
	configFile := filepath.Join(config.RootPath, "config", "config.toml")
	if _, err := os.Stat(configFile); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Error: go-switch is not initialized, please run 'goswitch init' first")
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Failed to check config file: %v\n", err)
		os.Exit(1)
	}

	if err := features.LoadConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	var cmd string
	args := os.Args
	if len(args) >= 2 {
		cmd = args[1]
	}

	switch cmd {
	case "help", "":
		PrintHelp()
	case "init":
		if err := features.InitGoSwitch(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize go-switch environment: %v\n", err)
			os.Exit(1)
		}
	case "listall":
		ensureInitialized()
		features.ListAll()
	case "install":
		ensureInitialized()
		searchVer := ""
		if len(args) >= 3 {
			searchVer = args[2]
		}
		if exists, create := helper.ExistsPath(config.GosPath); (exists || create) && searchVer != "" {
			features.Install(searchVer, string(config.SystemEnv), config.SystemArch, config.GosPath, config.TempUnzipPath)
		} else {
			panic("Please input the version you want to install")
		}
	case "switch":
		ensureInitialized()
		features.Switch()
	case "list":
		ensureInitialized()
		features.List()
	case "delete":
		ensureInitialized()
		features.Delete()
	case "env":
		ensureInitialized()
		features.Env()
	default:
		fmt.Println("Command not found")
	}
}

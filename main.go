package main

import (
	"fmt"
	"os"

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
	// 检查是否已初始化
	if config.RootPath == "" {
		fmt.Println("错误：go-switch 尚未初始化，请先运行 'goswitch init'")
		os.Exit(1)
	}
	features.LoadConfig()
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
		features.InitGoSwitch()
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

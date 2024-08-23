package main

import (
	"fmt"
	"os"

	"github.com/xulimeng/go-switch/config"
	"github.com/xulimeng/go-switch/features"
)

func PrintHelp() {
	fmt.Println("Usage: goswitch -cmd <command>")

	fmt.Println("\ngoswitch is the Go Version Manager")

	fmt.Println(`
Command:
	help	- Show this help message
	install	- Install to go version
	switch	- Choose go version
	list	- List all installed go versions
	listall - List all available go versions
	delete	- Delete go version
	env 	- Show goswitch environment
	`)
}

func init() {
	config.InitSystemVars()
	config.InitConfigFile()
	config.LoadConfig()
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
	case "listall":
		features.ListAll()
	case "install":
		searchVer := ""
		if len(args) >= 3 {
			searchVer = args[2]
		}
		if exists, create := config.ExistsPath(config.GosPath); (exists || create) && searchVer != "" {
			features.Install(searchVer, string(config.SystemEnv), config.SystemArch, config.GosPath, config.TempUnzipPath)
		} else {
			panic("Please input the version you want to install")
		}
	case "switch":
		features.Switch()
	case "list":
		features.List()
	case "delete":
		features.Delete()
	default:
		fmt.Println("Command not found")
	}
}

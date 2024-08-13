package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/xulimeng/go-switch/config"
	"github.com/xulimeng/go-switch/features"
	"github.com/xulimeng/go-switch/utils"
)

var (
	// 不同系统默认的 go 安装路径
	LinuxGoPath   = fmt.Sprintf("%s/", os.Getenv("HOME"))
	WindowsGoPath = `C:\\Users\\`
	MacGoPath     = fmt.Sprintf("%s/", os.Getenv("HOME"))

	SystemEnv     config.Env
	SystemArch    string
	RootPath      string
	GosPath       string
	TempUnzipPath string
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
	os := runtime.GOOS
	switch os {
	case "linux":
		SystemEnv = config.Linux
		RootPath = LinuxGoPath + config.GoSwitchDir
		GosPath = RootPath + "/" + config.SaveGoDir
		TempUnzipPath = GosPath + "/" + config.UnzipGoDir
	case "windows":
		SystemEnv = config.Windows
		userNameCurr, err := user.Current()
		if err != nil {
			panic(err)
		}
		RootPath = WindowsGoPath + userNameCurr.Username + "\\" + config.GoSwitchDir
		GosPath = RootPath + "\\" + config.SaveGoDir
		TempUnzipPath = GosPath + "\\" + config.UnzipGoDir
	case "darwin":
		SystemEnv = config.Mac
		RootPath = MacGoPath + config.GoSwitchDir
		GosPath = RootPath + "/" + config.SaveGoDir
		TempUnzipPath = GosPath + "/" + config.UnzipGoDir
	}
	SystemArch = runtime.GOARCH
}

func main() {

	fmt.Println("SystemEnv: ", SystemEnv)

	var cmd string
	args := os.Args
	if len(args) >= 2 {
		cmd = args[1]
	}

	if exists, create := utils.ExistsPath(RootPath); !exists && !create {
		panic("RootPath not exists")
	}

	// 初始化配置文
	config.InitConfigFile()
	fmt.Println("-----------Command: ", cmd)

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
		if exists, create := utils.ExistsPath(GosPath); (exists || create) && searchVer != "" {
			features.Install(searchVer, string(SystemEnv), SystemArch, GosPath, TempUnzipPath)
		} else {
			panic("Please input the version you want to install")
		}
	default:
		fmt.Println("Command not found")
	}

	// 更新配置文件
	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	if err := encoder.Encode(config.Conf); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(RootPath+"/config/config.toml", buffer.Bytes(), 0644); err != nil {
		panic(err)
	}

}

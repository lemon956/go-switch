package main

import (
	"fmt"
	"os"
	"os/user"
	"runtime"

	"github.com/xulimeng/go-switch/features"
)

type Env string

const (
	// Linux, Windows, Mac 系统环境类型
	Linux   Env = "linux"
	Windows Env = "windows"
	Mac     Env = "mac"
	// 不同系统默认的 go 安装路径
	LinuxGoPath   = "$HOME/"
	WindowsGoPath = `C:\\Users\\`
	MacGoPath     = "$HOME/"
	// go-switch 的文件夹名
	GoSwitchDir = ".go-switch"
	// 真正保存go 版本的文件夹名
	SaveGoDir = "gos"
)

var (
	SystemEnv Env
	RootPath  string
	GosPath   string
)

func PrintHelp() {
	fmt.Println("Usage: goswitch -cmd <command>")

	fmt.Println("\ngoswitch is the Go Version Manager")

	fmt.Println(`
Command:
	help	- Show this help message
	install	- Install to go version
	use		- Use go version
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
		SystemEnv = Linux
		RootPath = LinuxGoPath + GoSwitchDir
		GosPath = RootPath + "/" + SaveGoDir
	case "windows":
		SystemEnv = Windows
		userNameCurr, err := user.Current()
		if err != nil {
			panic(err)
		}
		RootPath = WindowsGoPath + userNameCurr.Username + "\\" + GoSwitchDir
		GosPath = RootPath + "\\" + SaveGoDir
	case "darwin":
		SystemEnv = Mac
		RootPath = MacGoPath + GoSwitchDir
		GosPath = RootPath + "/" + SaveGoDir
	}
}

// ExistsPath check if path exists
// Returns: bool, bool (exists, crate)
func ExistsPath(path string) (bool, bool) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return false, false
		}
		return false, true
	}
	return err == nil, false
}

func main() {

	fmt.Println("SystemEnv: ", SystemEnv)

	var cmd string
	args := os.Args
	if len(args) >= 2 {
		cmd = args[1]
	}

	fmt.Println("-----------Command: ", cmd)

	switch cmd {
	case "help", "":
		PrintHelp()
	case "listall":
		features.ListAll()
	case "install":
		features.Install()
	default:
		fmt.Println("Command not found")
	}

}

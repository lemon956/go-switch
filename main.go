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
	// go-switch 的文件夹名
	GoSwitchDir = ".go-switch"
	// 不同系统默认的 go 安装路径
	LinuxGoPath   = "$HOME/"
	WindowsGoPath = `C:\\Users\\`
	MacGoPath     = "$HOME/"
)

var (
	SystemEnv Env
	RootPath  string
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
	case "windows":
		SystemEnv = Windows
		userNameCurr, err := user.Current()
		if err != nil {
			panic(err)
		}
		RootPath = WindowsGoPath + userNameCurr.Username + "\\" + GoSwitchDir
	case "darwin":
		SystemEnv = Mac
		RootPath = MacGoPath + GoSwitchDir
	}
}

// ExistsPath check if path exists
// Returns: bool, bool (exists, isDir)
func ExistsPath(path string) (bool, bool) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, false
	}
	return err == nil, info.IsDir()
}

func main() {

	fmt.Println("SystemEnv: ", SystemEnv)

	var cmd string
	args := os.Args
	if len(args) > 2 {
		cmd = args[1]
	}

	fmt.Println("-----------Command: ", cmd)

	switch cmd {
	case "help", "":
		PrintHelp()
	case "listall":
		features.ListAll()
	}

}

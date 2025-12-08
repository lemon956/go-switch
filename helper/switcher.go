package helper

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xulimeng/go-switch/config"
)

var GlobalSwitcher Switcher

type Switcher interface {
	UpdateGoEnv(goRoot string)
	// 新的方法：使用软链接方式切换
	SwitchBySymlink(goVersion string) error
}

// createSymlink 创建或更新软链接
func createSymlink(source, target string) error {
	// 删除现有的软链接或文件
	if _, err := os.Lstat(target); err == nil {
		if err := os.Remove(target); err != nil {
			return fmt.Errorf("Failed to remove existing symlink target: %v", err)
		}
	}

	// 创建新的软链接
	if err := os.Symlink(source, target); err != nil {
		return fmt.Errorf("Failed to create symlink: %v", err)
	}

	return nil
}

// getGoSwitchBinPath 获取go-switch管理的Go二进制目录路径
func getGoSwitchBinPath() string {
	return filepath.Join(config.RootPath, "current", "bin")
}

// ensureGoSwitchBinInPath 确保go-switch的bin目录在PATH中
func ensureGoSwitchBinInPath() {
	goSwitchBin := getGoSwitchBinPath()

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(goSwitchBin), 0755); err != nil {
		fmt.Printf("Failed to create directory %s: %v\n", filepath.Dir(goSwitchBin), err)
		return
	}

	fmt.Printf("Please add the following directory to your PATH:\n%s\n", goSwitchBin)
	fmt.Println("You can add this line to your shell config file (e.g. ~/.bashrc or ~/.zshrc):")
	fmt.Printf("export PATH=\"%s:$PATH\"\n", goSwitchBin)
}

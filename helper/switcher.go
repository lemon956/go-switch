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
			return fmt.Errorf("无法删除现有目标: %v", err)
		}
	}

	// 创建新的软链接
	if err := os.Symlink(source, target); err != nil {
		return fmt.Errorf("无法创建软链接: %v", err)
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
		fmt.Printf("无法创建目录 %s: %v\n", filepath.Dir(goSwitchBin), err)
		return
	}

	fmt.Printf("请将以下目录添加到您的 PATH 环境变量中：\n%s\n", goSwitchBin)
	fmt.Println("您可以将以下命令添加到您的 shell 配置文件中（如 ~/.bashrc 或 ~/.zshrc）：")
	fmt.Printf("export PATH=\"%s:$PATH\"\n", goSwitchBin)
}

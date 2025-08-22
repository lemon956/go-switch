package features

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xulimeng/go-switch/config"
	"github.com/xulimeng/go-switch/helper"
)

// Env 显示go-switch环境信息
func Env() {
	fmt.Println("=== Go-Switch 环境信息 ===")
	fmt.Printf("Go-Switch 根目录: %s\n", config.RootPath)
	fmt.Printf("Go版本存储目录: %s\n", config.GosPath)
	fmt.Printf("系统环境: %s\n", config.SystemEnv)
	fmt.Printf("系统架构: %s\n", config.SystemArch)

	// 显示当前活跃的Go版本
	currentLinkPath := filepath.Join(config.RootPath, "current")
	if currentTarget, err := os.Readlink(currentLinkPath); err == nil {
		// 提取版本号
		versionName := filepath.Base(currentTarget)
		fmt.Printf("当前Go版本: %s\n", versionName)
		fmt.Printf("当前Go路径: %s\n", currentTarget)

		// 检查Go二进制是否可用
		goBinPath := filepath.Join(currentTarget, "bin", "go")
		if config.SystemEnv == config.Windows {
			goBinPath += ".exe"
		}

		if _, err := os.Stat(goBinPath); err == nil {
			fmt.Printf("Go二进制路径: %s\n", goBinPath)
		} else {
			fmt.Printf("警告：Go二进制不存在于 %s\n", goBinPath)
		}
	} else {
		fmt.Println("当前Go版本: 未设置")
		fmt.Println("提示：使用 'goswitch switch' 来选择一个Go版本")
	}

	// 显示已安装的Go版本
	fmt.Println("\n已安装的Go版本:")
	if len(config.Conf.LocalGos) > 0 {
		for _, goInfo := range config.Conf.LocalGos {
			fmt.Printf("  - %s (%s)\n", goInfo.Version, goInfo.Path)
		}
	} else {
		fmt.Println("  无已安装的版本")
		fmt.Println("  使用 'goswitch install <version>' 来安装Go版本")
	}

	// 显示PATH建议
	goSwitchBin := filepath.Join(config.RootPath, "current", "bin")
	fmt.Printf("\n建议的PATH配置:\n")
	fmt.Printf("export PATH=\"%s:$PATH\"\n", goSwitchBin)

	if config.SystemEnv != config.Windows {
		// Unix系统显示shell配置文件建议
		shell := helper.JudgeZshOrBash()
		var configFile string
		switch shell {
		case "zsh":
			configFile = "~/.zshrc"
		case "bash":
			configFile = "~/.bashrc"
			if config.SystemEnv == config.Mac {
				configFile = "~/.bash_profile"
			}
		}
		if configFile != "" {
			fmt.Printf("将上述命令添加到 %s 文件中\n", configFile)
		}
	}
}

// env_win.go
//go:build windows
// +build windows

package helper

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xulimeng/go-switch/config"
	"golang.org/x/sys/windows/registry"
)

type WinDowsSwitcher struct{}

func init() {
	GlobalSwitcher = &WinDowsSwitcher{}
}

const (
	HKEY_CURRENT_USER  = 0x80000001
	HKEY_LOCAL_MACHINE = 0x80000002
)

func (sw *WinDowsSwitcher) UpdateGoEnv(goRoot string) {
	err := setEnvVar(registry.CURRENT_USER, "GOROOT", goRoot)
	if err != nil {
		panic(err)
	}
	keyValue, err := getEnvVar(registry.CURRENT_USER, "PATH")
	if err != nil {
		panic(err)
	}
	err = setEnvVar(registry.CURRENT_USER, "PATH", fmt.Sprintf("%%GOROOT%%%sbin;%s", string(os.PathSeparator), keyValue))
	if err != nil {
		panic(err)
	}
}

func getEnvVar(key registry.Key, envKey string) (string, error) {
	k, err := registry.OpenKey(key, `Environment`, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()
	value, _, err := k.GetStringValue(envKey)
	return value, err
}

func setEnvVar(key registry.Key, envVar, value string) error {
	k, _, err := registry.CreateKey(key, `Environment`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	err = k.SetStringValue(envVar, value)
	if err != nil {
		return err
	}

	return nil
}

// SwitchBySymlink 使用软链接方式切换Go版本（Windows版本）
func (sw *WinDowsSwitcher) SwitchBySymlink(goVersion string) error {
	// 源目录：指定版本的Go安装目录
	sourceDir := filepath.Join(config.GosPath, goVersion)

	// 目标目录：go-switch管理的当前Go目录
	targetDir := filepath.Join(config.RootPath, "current")

	// 检查源目录是否存在
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return fmt.Errorf("Go版本 %s 不存在，请先安装", goVersion)
	}

	// 创建软链接（Windows上使用符号链接）
	if err := createSymlink(sourceDir, targetDir); err != nil {
		return fmt.Errorf("切换失败: %v", err)
	}

	// 在Windows上，设置GOROOT环境变量指向current目录
	if err := setEnvVar(registry.CURRENT_USER, "GOROOT", targetDir); err != nil {
		return fmt.Errorf("设置GOROOT失败: %v", err)
	}

	// 更新PATH环境变量
	binPath := filepath.Join(targetDir, "bin")
	currentPath, err := getEnvVar(registry.CURRENT_USER, "PATH")
	if err != nil {
		fmt.Printf("警告：无法获取当前PATH: %v\n", err)
	} else {
		newPath := fmt.Sprintf("%s;%s", binPath, currentPath)
		if err := setEnvVar(registry.CURRENT_USER, "PATH", newPath); err != nil {
			fmt.Printf("警告：无法更新PATH: %v\n", err)
		}
	}

	fmt.Printf("成功切换到 Go %s\n", goVersion)
	fmt.Printf("当前Go安装目录: %s\n", targetDir)
	fmt.Println("请重新启动命令行窗口以使环境变量生效")

	return nil
}

func UpdateGoEnvUnix() {
	fmt.Println("UpdateGoEnvUnix not in unix")
}

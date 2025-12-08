package features

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/lemon956/go-switch/config"
)

func Delete() {
	versions := []string{}
	if config.Conf.LocalGos == nil {
		config.Conf.LocalGos = make([]config.GosVersion, 0)
	}
	for _, vInfo := range config.Conf.LocalGos {
		versions = append(versions, vInfo.Version)
	}

	versions = append(versions, Exit)
	prompt := promptui.Select{
		Label: "Choose You Want Delete Version",
		Items: versions,
	}

	_, result, err := prompt.Run()
	if err != nil {
		panic(err)
	}
	if result == Exit {
		return
	}

	// 删除版本
	delIdx := -1
	for idx, value := range config.Conf.LocalGos {
		if value.Version == result {
			delIdx = idx
			break
		}
	}
	if delIdx == -1 {
		fmt.Println("Not Have This Version")
		return
	}

	// 删除文件系统中的Go版本目录
	versionPath := filepath.Join(config.GosPath, result)
	if err := os.RemoveAll(versionPath); err != nil {
		fmt.Printf("Failed to delete directory: %v\n", err)
		return
	}

	// 检查当前删除的版本是否是正在使用的版本
	currentLinkPath := filepath.Join(config.RootPath, "current")
	if currentTarget, err := os.Readlink(currentLinkPath); err == nil {
		if currentTarget == versionPath {
			// 如果删除的是当前使用的版本，移除软链接
			if err := os.Remove(currentLinkPath); err != nil {
				fmt.Printf("Warning: failed to remove the currently active Go version: %v\n", err)
			} else {
				fmt.Println("The currently active Go version has been removed, please switch to another version")
			}
			config.Conf.GoRoot = ""
		}
	}

	// 从配置中移除该版本
	config.Conf.LocalGos = append(config.Conf.LocalGos[:delIdx], config.Conf.LocalGos[delIdx+1:]...)
	config.Conf.SaveConfig()

	fmt.Printf("Deleted Go %s successfully\n", result)
}

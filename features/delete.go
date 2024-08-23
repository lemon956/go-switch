package features

import (
	"fmt"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/xulimeng/go-switch/config"
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

	// 清理环境变量
	goRootPath := filepath.Join(config.GosPath, result)
	GlobalSwitcher.UpdateGoEnv(goRootPath)
	config.Conf.GoRoot = goRootPath
	config.Conf.LocalGos = append(config.Conf.LocalGos[:delIdx], config.Conf.LocalGos[delIdx+1:]...)
	config.Conf.SaveConfig()
}

package features

import (
	"fmt"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/xulimeng/go-switch/config"
	"github.com/xulimeng/go-switch/helper"
)

const Exit = "exit"

func Switch() {
	versions := []string{}
	if config.Conf.LocalGos == nil {
		config.Conf.LocalGos = make([]config.GosVersion, 0)
	}
	for _, vInfo := range config.Conf.LocalGos {
		versions = append(versions, vInfo.Version)
	}

	if len(versions) == 0 {
		fmt.Println("No installed Go versions found, please run 'goswitch install' first")
		return
	}

	versions = append(versions, Exit)
	prompt := promptui.Select{
		Label: "Choose You Want Switch Version",
		Items: versions,
	}

	_, result, err := prompt.Run()
	if err != nil {
		panic(err)
	}
	if result == Exit {
		return
	}

	// 使用新的软链接方式切换
	if err := helper.GlobalSwitcher.SwitchBySymlink(result); err != nil {
		fmt.Printf("Switch failed: %v\n", err)
		return
	}

	// 更新配置文件中的当前Go版本
	config.Conf.GoRoot = filepath.Join(config.RootPath, "current")
	config.Conf.SaveConfig()

	fmt.Printf("Switched to Go %s successfully\n", result)
}

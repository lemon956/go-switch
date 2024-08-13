package features

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/xulimeng/go-switch/config"
	"github.com/xulimeng/go-switch/utils"
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

	// TODO: switch version
}

func UpdateGoEnv(goRoot string) {
	// set GOROOT
	if config.SystemEnv == config.Linux || config.SystemEnv == config.Mac {
		sh := utils.JudgeZshOrBash()
		switch sh {
		case "zsh":

			break
		case "bash":
			break
		default:
			fmt.Println("Not support shell")
		}
	}
	err := os.Setenv("GOROOT", goRoot)
	if err != nil {
		panic(err)
	}
	// set PATH
}

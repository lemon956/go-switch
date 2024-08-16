// env_win.go
//go:build windows
// +build windows

package features

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows/registry"
)

func UpdateGoEnvWin() {
	vallue := os.Getenv("PATH")
	fmt.Println("os.Getenv ", vallue)

	k, _, err := registry.CreateKey(registry.CURRENT_USER, `Environment`, registry.SET_VALUE)
	if err != nil {
		panic(err)
	}
	defer k.Close()
	value, _, err := k.GetStringValue("PATH")
	if err != nil {
		panic(err)
	}
	fmt.Println("registry.GetStringValue", value)
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

func UpdateGoEnvUnix() {
	fmt.Println("UpdateGoEnvUnix not in unix")
}

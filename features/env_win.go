// env_win.go
//go:build windows
// +build windows

package features

import (
	"fmt"
	"os"

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

func UpdateGoEnvUnix() {
	fmt.Println("UpdateGoEnvUnix not in unix")
}

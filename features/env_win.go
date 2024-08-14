// env_win.go
//go:build windows
// +build windows

package features

import (
	"golang.org/x/sys/windows/registry"
)

func UpdateGoEnvWin() {
	setEnvVar(registry.CURRENT_USER, "", "")
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

package features

import (
	"os"
	"testing"

	"github.com/xulimeng/go-switch/config"
	"github.com/xulimeng/go-switch/helper"
)

func TestMain(m *testing.M) {
	InitSystemVars()
	if exists, create := helper.ExistsPath(config.RootPath); !exists && !create {
		panic("RootPath not exists")
	}
	if err := InitConfigFile(); err != nil {
		panic(err)
	}
	if err := LoadConfig(); err != nil {
		panic(err)
	}

	code := m.Run()
	os.Exit(code)
}

func TestSwitch(t *testing.T) {
	Switch()
}

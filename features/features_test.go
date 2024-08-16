package features

import (
	"os"
	"testing"

	"github.com/xulimeng/go-switch/config"
)

func TestMain(m *testing.M) {
	config.InitSystemVars()
	if exists, create := config.ExistsPath(config.RootPath); !exists && !create {
		panic("RootPath not exists")
	}
	config.InitConfigFile()
	config.LoadConfig()

	code := m.Run()
	os.Exit(code)
}

func TestSwitch(t *testing.T) {
	Switch()
}

func Test(t *testing.T) {
	rsp := JudgeZshOrBash()
	t.Log(rsp)
}

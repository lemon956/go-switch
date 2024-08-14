package features

import (
	"os"
	"testing"

	"github.com/xulimeng/go-switch/config"
	"github.com/xulimeng/go-switch/utils"
)

func TestMain(m *testing.M) {
	config.InitSystemVars()
	if exists, create := utils.ExistsPath(config.RootPath); !exists && !create {
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

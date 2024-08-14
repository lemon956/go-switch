package config

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	RootPath = "/home/hellotalk/.go-switch"
	LoadConfig()
	code := m.Run()
	os.Exit(code)
}

func TestLoadConfig(t *testing.T) {
	t.Logf("%+v", Conf)
}

package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	GoSwitchPath string       `toml:"go_switch_path"`
	Init         bool         `toml:"init"`
	LocalGos     []GosVersion `toml:"local_gos"`

	LocalGolangCILints  []GolangCILintVersion `toml:"local_golangci_lints"`
	CurrentGolangCILint GolangCILintCurrent   `toml:"current_golangci_lint"`

	// 当前生效的 golang 环境变量
	GoPath string `toml:"go_path"`
	GoRoot string `toml:"go_root"`
}

type GosVersion struct {
	Version string `toml:"version"`
	Path    string `toml:"path"`
}

type GolangCILintVersion struct {
	Major      string `toml:"major"`
	Version    string `toml:"version"`
	Path       string `toml:"path"`
	BinaryPath string `toml:"binary_path"`
}

type GolangCILintCurrent struct {
	Major      string `toml:"major"`
	Version    string `toml:"version"`
	BinaryPath string `toml:"binary_path"`
}

var Conf *Config

func (c *Config) SaveConfig() {
	// 更新配置文件
	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	if err := encoder.Encode(Conf); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(fmt.Sprintf("%s%s%s", filepath.Join(RootPath, "config"), string(os.PathSeparator), "config.toml"), buffer.Bytes(), 0644); err != nil {
		panic(err)
	}
}

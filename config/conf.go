package config

type Config struct {
	GoSwitchPath string       `toml:"go_switch_path"`
	LocalGos     []GosVersion `toml:"local_gos"`

	// 当前生效的 golang 环境变量
	GoPath string `toml:"go_path"`
	GoRoot string `toml:"go_root"`
}

type GosVersion struct {
	Version string `toml:"version"`
	Path    string `toml:"path"`
}

var Conf *Config

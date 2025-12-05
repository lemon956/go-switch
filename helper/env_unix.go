// env_unix.go
//go:build !windows
// +build !windows

package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xulimeng/go-switch/config"
)

type UnixSwitcher struct{}

func init() {
	GlobalSwitcher = &UnixSwitcher{}
}

// getShellEnvFile returns the environment file path for the given shell
func getShellEnvFile(shell string) string {
	if shell == "fish" {
		return filepath.Join(config.GoEnvFilePath, "system.fish")
	}
	return filepath.Join(config.GoEnvFilePath, "system")
}

// getShellConfigFile returns the shell config file path for the given shell
func getShellConfigFile(shell string) string {
	home := os.Getenv("HOME")
	switch shell {
	case "zsh":
		return home + "/.zshrc"
	case "bash":
		if config.SystemEnv == config.Mac {
			return home + "/.bash_profile"
		}
		return home + "/.bashrc"
	case "fish":
		return home + "/.config/fish/config.fish"
	}
	return ""
}

// generateGoRootCmd generates the GOROOT command for the given shell
func generateGoRootCmd(shell, goRoot string) string {
	if shell == "fish" {
		return fmt.Sprintf("set -gx GOROOT %s", goRoot)
	}
	return fmt.Sprintf("export GOROOT=%s", goRoot)
}

// generatePathCmd generates the PATH command for the given shell
func generatePathCmd(shell string) string {
	// Prefer to keep GOPATH/bin first if GOPATH is configured
	if shell == "fish" {
		if config.Conf != nil && config.Conf.GoPath != "" {
			// Use GOPATH and GOROOT bins
			return "set -gx PATH $GOPATH/bin $GOROOT/bin $PATH"
		}
		return "set -gx PATH $GOROOT/bin $PATH"
	}

	if config.Conf != nil && config.Conf.GoPath != "" {
		return "export PATH=$GOPATH/bin:$GOROOT/bin:$PATH"
	}
	return "export PATH=$GOROOT/bin:$PATH"
}

// generateSourceCmd generates the source command for the given shell
func generateSourceCmd(shell, envFile string) string {
	if shell == "fish" {
		// In fish, check that the environment file exists before sourcing it
		// Use a single line so it can be safely appended to config.fish
		return fmt.Sprintf("if test -f %s; source %s; end", envFile, envFile)
	}
	return fmt.Sprintf("source %s", envFile)
}

// UpdateGoEnv updates Unix system environment variables
func (sw *UnixSwitcher) UpdateGoEnv(goRoot string) {
	sh := DetectShell()

	// Get shell-specific paths
	goEnvFilePath := getShellEnvFile(sh)
	configFile := getShellConfigFile(sh)

	if configFile == "" {
		fmt.Println("Current shell is not supported for automatic config update.")
		fmt.Println("Supported shells: bash, zsh, fish")
		fmt.Println("Please manually add the go-switch environment to your shell config.")
		return
	}

	// Ensure environment directory exists
	envDir := filepath.Dir(goEnvFilePath)
	if _, err := os.Stat(envDir); os.IsNotExist(err) {
		if err := os.MkdirAll(envDir, 0755); err != nil {
			fmt.Printf("Failed to create environment directory: %v\n", err)
			return
		}
	}

	// For fish, ensure config directory exists
	if sh == "fish" {
		fishConfigDir := filepath.Dir(configFile)
		if _, err := os.Stat(fishConfigDir); os.IsNotExist(err) {
			if err := os.MkdirAll(fishConfigDir, 0755); err != nil {
				fmt.Printf("Failed to create fish config directory: %v\n", err)
				return
			}
		}
	}

	// Generate shell-specific commands
	goRootCmd := generateGoRootCmd(sh, goRoot)
	pathCmd := generatePathCmd(sh)

	// Write to environment file
	if config.GoEnvFilePath != "" {
		addEnvironmentVariable(goEnvFilePath, goRootCmd, sh)
		addEnvironmentVariable(goEnvFilePath, pathCmd, sh)
	}

	// Add source command to shell config; addEnvironmentVariable is idempotent per variable type
	if configFile != "" && goEnvFilePath != "" {
		sourceCmd := generateSourceCmd(sh, goEnvFilePath)
		addEnvironmentVariable(configFile, sourceCmd, sh)
	}

	fmt.Println("go-switch environment files have been updated:")
	fmt.Printf("  - Environment file: %s\n", goEnvFilePath)
	fmt.Printf("  - Shell config: %s\n", configFile)
	fmt.Println("Hint: If this is your first time using goswitch in this shell, run the following command:")
	fmt.Printf("  source %s\n", configFile)
}

// getEnvVarPrefix returns the prefix patterns for environment variables based on shell type
func getEnvVarPrefix(shell, varName string) (setPrefix, checkPrefix string) {
	if shell == "fish" {
		setPrefix = fmt.Sprintf("set -gx %s ", varName)
		checkPrefix = fmt.Sprintf("set -gx %s", varName)
	} else {
		setPrefix = fmt.Sprintf("export %s=", varName)
		checkPrefix = fmt.Sprintf("export %s=", varName)
	}
	return
}

// addEnvironmentVariable adds or updates an environment variable in the config file
func addEnvironmentVariable(configFile, line, shell string) {
	// Read file content
	content, err := os.ReadFile(configFile)
	if err != nil {
		// If file doesn't exist, create it and write
		if os.IsNotExist(err) {
			// Ensure parent directory exists
			dir := filepath.Dir(configFile)
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Printf("Failed to create directory %s: %v\n", dir, err)
				return
			}
			if err := os.WriteFile(configFile, []byte(line+"\n"), 0644); err != nil {
				fmt.Printf("Failed to create file %s: %v\n", configFile, err)
			} else {
				fmt.Printf("Created file and added '%s' to %s\n", line, configFile)
			}
			return
		}
		fmt.Printf("Failed to read %s: %v\n", configFile, err)
		return
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	lineAdded := false

	// Detect what type of variable this is
	_, goRootCheck := getEnvVarPrefix(shell, "GOROOT")
	_, pathCheck := getEnvVarPrefix(shell, "PATH")

	isGoRoot := strings.HasPrefix(line, goRootCheck)
	isPath := strings.HasPrefix(line, pathCheck)

	for _, existingLine := range lines {
		existingLine = strings.TrimSpace(existingLine)

		// If the exact same line already exists, keep it and avoid adding duplicates
		if existingLine != "" && strings.TrimSpace(line) == existingLine {
			newLines = append(newLines, existingLine)
			lineAdded = true
			continue
		}

		if isGoRoot && strings.HasPrefix(existingLine, goRootCheck) {
			// GOROOT: overwrite
			newLines = append(newLines, line)
			lineAdded = true
			fmt.Printf("Replaced GOROOT in %s with '%s'\n", configFile, line)
		} else if isPath && strings.HasPrefix(existingLine, pathCheck) {
			// PATH: keep if same, replace if different
			if existingLine == line {
				newLines = append(newLines, existingLine)
				lineAdded = true
				fmt.Printf("PATH already exists in %s: '%s'\n", configFile, line)
			} else {
				newLines = append(newLines, line)
				lineAdded = true
				fmt.Printf("Replaced PATH in %s with '%s'\n", configFile, line)
			}
		} else if existingLine != "" {
			// Keep other lines
			newLines = append(newLines, existingLine)
		}
	}

	// If no matching line found, add new line
	if !lineAdded {
		newLines = append(newLines, line)
		fmt.Printf("Added '%s' to %s\n", line, configFile)
	}

	// Write back to file
	newContent := strings.Join(newLines, "\n") + "\n"
	if err := os.WriteFile(configFile, []byte(newContent), 0644); err != nil {
		fmt.Printf("Failed to write to %s: %v\n", configFile, err)
	}
}

// SwitchBySymlink 使用软链接方式切换Go版本
func (sw *UnixSwitcher) SwitchBySymlink(goVersion string) error {
	// 源目录：指定版本的Go安装目录
	sourceDir := filepath.Join(config.GosPath, goVersion)

	// 目标目录：go-switch管理的当前Go目录
	targetDir := filepath.Join(config.RootPath, "current")

	// 检查源目录是否存在
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return fmt.Errorf("go version %s does not exist, please install it first", goVersion)
	}

	// 创建软链接
	if err := createSymlink(sourceDir, targetDir); err != nil {
		return fmt.Errorf("switch failed: %v", err)
	}

	// 更新 go-switch 环境变量配置（GOROOT / PATH）
	sw.UpdateGoEnv(targetDir)

	fmt.Printf("Switched to Go %s successfully\n", goVersion)
	fmt.Printf("Current Go install path: %s\n", targetDir)

	return nil
}

func UpdateGoEnvWin() {
	fmt.Println("UpdateGoEnvWin not in windows")
}

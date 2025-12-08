package features

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lemon956/go-switch/config"
	"github.com/lemon956/go-switch/helper"
)

// Env displays go-switch environment information
func Env() {
	fmt.Println("=== Go-Switch Environment ===")
	fmt.Printf("Go-Switch root: %s\n", config.RootPath)
	fmt.Printf("Go versions directory: %s\n", config.GosPath)
	fmt.Printf("System OS: %s\n", config.SystemEnv)
	fmt.Printf("System arch: %s\n", config.SystemArch)

	// Display current active Go version
	currentLinkPath := filepath.Join(config.RootPath, "current")
	if currentTarget, err := os.Readlink(currentLinkPath); err == nil {
		// Extract version number
		versionName := filepath.Base(currentTarget)
		fmt.Printf("Current Go version: %s\n", versionName)
		fmt.Printf("Current Go path: %s\n", currentTarget)

		// Check if Go binary is available
		goBinPath := filepath.Join(currentTarget, "bin", "go")
		if config.SystemEnv == config.Windows {
			goBinPath += ".exe"
		}

		if _, err := os.Stat(goBinPath); err == nil {
			fmt.Printf("Go binary path: %s\n", goBinPath)
		} else {
			fmt.Printf("Warning: Go binary does not exist at %s\n", goBinPath)
		}
	} else {
		fmt.Println("Current Go version: not set")
		fmt.Println("Hint: use 'goswitch switch' to select a Go version")
	}

	// Display installed Go versions
	fmt.Println("\nInstalled Go versions:")
	if len(config.Conf.LocalGos) > 0 {
		for _, goInfo := range config.Conf.LocalGos {
			fmt.Printf("  - %s (%s)\n", goInfo.Version, goInfo.Path)
		}
	} else {
		fmt.Println("  No versions installed")
		fmt.Println("  Use 'goswitch install <version>' to install a Go version")
	}

	if config.SystemEnv != config.Windows {
		// Unix systems: show shell-specific environment loading instructions
		shell := helper.DetectShell()
		var goEnvSystem, configFile, sourceCmd string

		switch shell {
		case "fish":
			goEnvSystem = filepath.Join(config.GoEnvFilePath, "system.fish")
			configFile = "~/.config/fish/config.fish"
			sourceCmd = fmt.Sprintf("source %s", goEnvSystem)
		case "zsh":
			goEnvSystem = filepath.Join(config.GoEnvFilePath, "system")
			configFile = "~/.zshrc"
			sourceCmd = fmt.Sprintf("source %s", goEnvSystem)
		case "bash":
			goEnvSystem = filepath.Join(config.GoEnvFilePath, "system")
			if config.SystemEnv == config.Mac {
				configFile = "~/.bash_profile"
			} else {
				configFile = "~/.bashrc"
			}
			sourceCmd = fmt.Sprintf("source %s", goEnvSystem)
		default:
			goEnvSystem = filepath.Join(config.GoEnvFilePath, "system")
			configFile = ""
			sourceCmd = fmt.Sprintf("source %s", goEnvSystem)
		}

		fmt.Printf("\nEnvironment file (maintained by goswitch):\n  %s\n", goEnvSystem)
		fmt.Println("\nAdd the following line to your shell config (only once):")
		fmt.Printf("  %s\n", sourceCmd)

		if configFile != "" {
			fmt.Printf("Add the command above into %s\n", configFile)
		}

		// Show all supported shells
		fmt.Println("\nSupported shells: bash, zsh, fish")
		fmt.Printf("Detected shell: %s\n", shell)
	} else {
		// Windows: show generic environment file path
		goEnvSystem := filepath.Join(config.GoEnvFilePath, "system")
		fmt.Printf("\nEnvironment file (maintained by goswitch):\n  %s\n", goEnvSystem)
	}
}

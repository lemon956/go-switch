// env_unix.go
//go:build !windows
// +build !windows

package features

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/xulimeng/go-switch/config"
)

type UnixSwitcher struct{}

func init() {
	GlobalSwitcher = &UnixSwitcher{}
}

// UpdateGoEnvUnix 更新 Unix 系统的环境变量
func (sw *UnixSwitcher) UpdateGoEnv(goRoot string) {
	// set GOROOT
	sh := JudgeZshOrBash()
	goRootCmd := fmt.Sprintf("export GOROOT=%s", goRoot)
	pathCmd := "export PATH=$GOROOT/bin:$PATH"
	goEnvFilePath := fmt.Sprintf("%s%s%s", config.GoEnvFilePath, string(os.PathSeparator), "system")
	if config.GoEnvFilePath != "" {
		addEnvironmentVariable(goEnvFilePath, goRootCmd)
		addEnvironmentVariable(goEnvFilePath, pathCmd)
	}
	var configFile string
	switch sh {
	case "zsh":
		configFile = os.Getenv("HOME") + "/.zshrc"
	case "bash":
		configFile = os.Getenv("HOME") + "/.bashrc"
		if config.SystemEnv == config.Mac {
			configFile = os.Getenv("HOME") + "/.bash_profile"
		}
	default:
		fmt.Println("Not support shell")
		return
	}
	if !config.Conf.Init && configFile != "" && goEnvFilePath != "" {
		addEnvironmentVariable(configFile, fmt.Sprintf("source %s", goEnvFilePath))
		config.Conf.Init = true
		config.Conf.SaveConfig()
	}
	// if err := reloadZshCOnfig(sh, configFile); err != nil {
	// 	fmt.Printf("Failed to reload %s config: %v\n", sh, err)
	// 	panic(err)
	// }
	// if !config.Conf.Init && configFile != "" && goEnvFilePath != "" {
	// 	config.Conf.Init = true
	// 	config.Conf.SaveConfig()
	// }
	fmt.Println("Please execute the following command: ")
	fmt.Println("source " + configFile)
}

// addEnvironmentVariable 添加环境变量
func addEnvironmentVariable(configFile, line string) {
	// 读取文件内容
	content, err := os.ReadFile(configFile)
	if err != nil {
		// 如果文件不存在，创建文件并写入
		if os.IsNotExist(err) {
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

	// 判断是GOROOT、PATH还是GOPATH
	isGoRoot := strings.HasPrefix(line, "export GOROOT=")
	isPath := strings.HasPrefix(line, "export PATH=")

	for _, existingLine := range lines {
		existingLine = strings.TrimSpace(existingLine)

		if isGoRoot && strings.HasPrefix(existingLine, "export GOROOT=") {
			// GOROOT：覆盖写入
			newLines = append(newLines, line)
			lineAdded = true
			fmt.Printf("Replaced GOROOT in %s with '%s'\n", configFile, line)
		} else if isPath && strings.HasPrefix(existingLine, "export PATH=") {
			// PATH：如果整行相同则保留，如果不同则替换
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
			// 保留其他行
			newLines = append(newLines, existingLine)
		}
	}

	// 如果没有找到对应的行，则添加新行
	if !lineAdded {
		newLines = append(newLines, line)
		fmt.Printf("Added '%s' to %s\n", line, configFile)
	}

	// 写回文件
	newContent := strings.Join(newLines, "\n") + "\n"
	if err := os.WriteFile(configFile, []byte(newContent), 0644); err != nil {
		fmt.Printf("Failed to write to %s: %v\n", configFile, err)
	}
}

func reloadZshCOnfig(shCmd string, shPath string) error {
	cmd := exec.Command(shCmd, "-c", fmt.Sprintf("source %s", shPath))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// JudgeZshOrBash 判断当前 shell 类型
func JudgeZshOrBash() string {
	// 获取 SHELL 环境变量
	shell := os.Getenv("SHELL")
	if shell == "" {
		fmt.Println("SHELL environment variable is not set")
		return ""
	}

	currentShell := ""
	shellSplit := strings.Split(shell, "/")
	if len(shellSplit) > 0 {
		currentShell = shellSplit[len(shellSplit)-1]
	}
	// 根据 shell 类型执行不同操作
	if strings.Contains(currentShell, "zsh") {
		return "zsh"
	} else if strings.Contains(currentShell, "bash") {
		return "bash"
	}
	return ""
}

func UpdateGoEnvWin() {
	fmt.Println("UpdateGoEnvWin not in windows")
}

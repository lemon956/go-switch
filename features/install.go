package features

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/xulimeng/go-switch/config"
	"github.com/xulimeng/go-switch/models"
	"github.com/xulimeng/go-switch/utils"
)

// Install 安装 Go 版本
func Install(searchVer string, system string, arch string, savePath string, unzipGoPath string) {
	// 获取 Go 版本信息
	resp, err := http.Get(models.GoVersionsURL)
	if err != nil {
		fmt.Println("Error fetching Go versions:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// 解析 JSON 数据
	var versions []models.GoVersion
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		fmt.Println("Connect Golang  Failed:", err)
		os.Exit(1)
	}

	// 打印每个版本的下载链接
	for _, version := range versions {
		if version.Version == searchVer {
			for _, file := range version.Files {
				if system == file.OS && arch == file.Arch {
					fmt.Printf("OS: %s\n", file.OS)
					fmt.Printf("Arch: %s\n", file.Arch)
					fmt.Printf("Filename: %s\n", file.Filename)
					fmt.Printf("Size: %d\n", file.Size)
					fmt.Printf("Kind: %s\n", file.Kind)
					fmt.Printf("Sha256: %s\n", file.Sha256)
					fmt.Println("Download URL: https://golang.org/dl/" + file.Filename)
					filePathName := savePath + "/" + file.Filename
					err := utils.DownloadFile(models.GoBaseURL+file.Filename, filePathName)
					if err != nil {
						panic(fmt.Sprintf("Download %s failed: %v", file.Filename, err))
					}
					err = utils.UnTarGz(filePathName, savePath)
					if err != nil {
						panic(fmt.Sprintf("UnTarGz %s failed: %v", file.Filename, err))
					}
					err = utils.RenameDir(unzipGoPath, version.Version)
					if err != nil {
						panic(fmt.Sprintf("RenameDir %s failed: %v", file.Filename, err))
					}
					err = os.Remove(filePathName)
					if err != nil {
						panic(fmt.Sprintf("Remove %s failed: %v", file.Filename, err))
					}
					afterRenamePath := config.GosPath + "/" + version.Version
					if config.SystemEnv == config.Windows {
						afterRenamePath = config.GosPath + "\\" + version.Version
					}
					config.Conf.LocalGos = append(config.Conf.LocalGos, afterRenamePath)
				}
			}
		}
	}
}

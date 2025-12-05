package features

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/xulimeng/go-switch/config"
	"github.com/xulimeng/go-switch/helper"
	"github.com/xulimeng/go-switch/models"
)

// Install 安装 Go 版本
func Install(searchVer string, system string, arch string, savePath string, unzipGoPath string) {
	// Create HTTP client with timeout to avoid hanging
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Fetch Go version metadata
	resp, err := client.Get(models.GoVersionsURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching Go versions from %s: %v\n", models.GoVersionsURL, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Error: received non-200 response code from Go versions API: %d %s\n", resp.StatusCode, resp.Status)
		os.Exit(1)
	}

	// 解析 JSON 数据
	var versions []models.GoVersion
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decode Go versions JSON: %v\n", err)
		os.Exit(1)
	}

	// 打印每个版本的下载链接
	for _, version := range versions {
		if version.Version == searchVer {
			for _, file := range version.Files {
				if system == file.OS && arch == file.Arch && file.Kind == "archive" {
					fmt.Printf("OS: %s\n", file.OS)
					fmt.Printf("Arch: %s\n", file.Arch)
					fmt.Printf("Filename: %s\n", file.Filename)
					fmt.Printf("Size: %d\n", file.Size)
					fmt.Printf("Kind: %s\n", file.Kind)
					fmt.Printf("Sha256: %s\n", file.Sha256)
					fmt.Println("Download URL: https://golang.org/dl/" + file.Filename)
					filePathName := savePath + string(os.PathSeparator) + file.Filename
					if exists, _ := helper.FileExists(filePathName); exists {
						// 删掉已存在的文件
						err := os.Remove(filePathName)
						if err != nil {
							panic(fmt.Sprintf("Remove %s failed: %v", file.Filename, err))
						}
					}
					err := helper.DownloadFile(models.GoBaseURL+file.Filename, filePathName)
					if err != nil {
						panic(fmt.Sprintf("Download %s failed: %v", file.Filename, err))
					}
					err = helper.Decompress(filePathName, savePath)
					if err != nil {
						panic(fmt.Sprintf("UnTarGz %s failed: %v", file.Filename, err))
					}
					err = helper.RenameDir(unzipGoPath, version.Version)
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

					if err := helper.GlobalSetPermissions.SetPermissions(afterRenamePath); err != nil {
						panic(fmt.Sprintf("SetPermissions %s failed: %v", file.Filename, err))
					}

					config.Conf.LocalGos = append(config.Conf.LocalGos, config.GosVersion{
						Version: version.Version,
						Path:    afterRenamePath,
					})
					config.Conf.SaveConfig()

					fmt.Printf("Go %s has been installed successfully!\n", version.Version)
					fmt.Println("Use 'goswitch switch' to switch to this version")
					fmt.Println("Use 'goswitch env' to inspect the environment")
					return
				}
			}
		}
	}
}

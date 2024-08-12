package features

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/xulimeng/go-switch/models"
)

func Install() {
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
		fmt.Printf("Go Version: %s\n", version.Version)
		for _, file := range version.Files {
			fmt.Printf("  OS: %s, Arch: %s, Download: %s\n", file.OS, file.Arch, file.Download)
		}
	}
}

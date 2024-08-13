package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func DownloadFile(url, filepath string) error {
	// 创建 HTTP 请求
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查 HTTP 响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server not response 200 code: %d %s", resp.StatusCode, resp.Status)
	}

	// 创建文件
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 获取内容长度
	contentLength := resp.ContentLength

	// 创建进度条
	bar := progressbar.NewOptions64(contentLength,
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("Downloading..."),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
	)

	// 创建一个多写入器，同时写入文件和进度条
	writer := io.MultiWriter(out, bar)

	// 将响应主体拷贝到文件
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// UnTarGz 解压 tar.gz 文件到指定目录
func UnTarGz(src, dest string) error {
	// 打开 tar.gz 文件
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建 gzip.Reader
	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	// 创建 tar.Reader
	tr := tar.NewReader(gzr)

	// 解压文件
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// 构建文件路径
		target := filepath.Join(dest, header.Name)

		// 检查文件类型
		switch header.Typeflag {
		case tar.TypeDir:
			// 创建目录
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			// 创建文件
			if err := os.MkdirAll(filepath.Dir(target), os.FileMode(header.Mode)); err != nil {
				return err
			}
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		default:
			fmt.Printf("无法识别的文件类型: %v\n", header.Typeflag)
		}
	}

	return nil
}

// RenameDir 重命名解压后的文件夹
func RenameDir(src, newName string) error {
	parentDir := filepath.Dir(src)
	newPath := filepath.Join(parentDir, newName)

	return os.Rename(src, newPath)
}

// ExistsPath check if path exists
// Returns: bool, bool (exists, crate)
func ExistsPath(path string) (bool, bool) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return false, false
		}
		return false, true
	}
	return err == nil, false
}

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

//go:build windows
// +build windows

// helper_windows.go
package config

import (
	"os"

	"golang.org/x/sys/windows"
)

func SetPermissionsWindows(path string) error {
	// 打开文件
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	// 获取文件句柄
	fileHandle := windows.Handle(file.Fd())

	// 获取当前文件的安全描述符
	sd, err := windows.GetSecurityInfo(fileHandle, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		return err
	}

	// 创建一个新的 ACL
	dacl, _, err := sd.DACL()
	if err != nil {
		return err
	}

	dacl, err = windows.ACLFromEntries([]windows.EXPLICIT_ACCESS{
		{
			AccessPermissions: windows.MAXIMUM_ALLOWED,
			AccessMode:        windows.GRANT_ACCESS,
			Inheritance:       windows.NO_INHERITANCE,
		},
	}, dacl)
	if err != nil {
		return err
	}

	// 设置新的 DACL 到文件的安全描述符中
	err = windows.SetSecurityInfo(fileHandle, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION, nil, nil, dacl, nil)
	if err != nil {
		return err
	}

	// 这里可以添加更多的 Windows 特定的权限设置代码

	return nil
}

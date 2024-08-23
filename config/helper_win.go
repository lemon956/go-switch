//go:build windows
// +build windows

// helper_windows.go
package config

import (
	"syscall"
)

type WindowsSetPermissions struct{}

func init() {
	GlobalSetPermissions = &WindowsSetPermissions{}
}

func (sp *WindowsSetPermissions) SetPermissions(path string) error {
	// // 打开文件
	// file, err := os.OpenFile(path, os.O_RDWR, 0)
	// if err != nil {
	// 	return err
	// }
	// defer file.Close()

	// // 获取文件句柄
	// fileHandle := windows.Handle(file.Fd())

	// // 获取当前文件的安全描述符
	// sd, err := windows.GetSecurityInfo(fileHandle, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	// if err != nil {
	// 	return err
	// }

	// // 创建一个新的 ACL
	// dacl, _, err := sd.DACL()
	// if err != nil {
	// 	return err
	// }

	// dacl, err = windows.ACLFromEntries([]windows.EXPLICIT_ACCESS{
	// 	{
	// 		AccessPermissions: windows.MAXIMUM_ALLOWED,
	// 		AccessMode:        windows.GRANT_ACCESS,
	// 		Inheritance:       windows.NO_INHERITANCE,
	// 	},
	// }, dacl)
	// if err != nil {
	// 	return err
	// }

	// // 设置新的 DACL 到文件的安全描述符中
	// err = windows.SetSecurityInfo(fileHandle, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION, nil, nil, dacl, nil)
	// if err != nil {
	// 	return err
	// }

	// // 这里可以添加更多的 Windows 特定的权限设置代码

	return nil
}

// 设置文件夹为隐藏属性
func (sp *WindowsSetPermissions) SetHiddenAttribute(path string) error {
	// 将路径转换为Windows API需要的UTF-16格式
	p, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	// 调用Windows API设置文件属性
	attrs, err := syscall.GetFileAttributes(p)
	if err != nil {
		return err
	}

	// 添加隐藏属性
	attrs |= syscall.FILE_ATTRIBUTE_HIDDEN

	// 设置新的文件属性
	return syscall.SetFileAttributes(p, attrs)
}

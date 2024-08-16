package config

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

func SetPermissionsUnix(path string) error {
	// 获取当前登录用户
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	// 解析 UID 和 GID
	uidStr := currentUser.Uid
	uid, _ := strconv.Atoi(uidStr)
	gidStr := currentUser.Gid
	gid, _ := strconv.Atoi(gidStr)
	// 使用 filepath.Walk 递归遍历目录树
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 改变文件或目录的权限
		err = os.Chmod(path, 0755)
		if err != nil {
			return err
		}

		// 改变文件或目录的所有权
		err = os.Chown(path, uid, gid)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

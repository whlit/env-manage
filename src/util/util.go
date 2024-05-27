package util

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// 获取根目录 获取失败则直接退出程序
func GetRootDir() (string) {
	exePath, err := os.Executable()
	if err != nil {
        fmt.Println("获取根目录失败")
		os.Exit(1)
	}
	return filepath.Dir(exePath)
}

// 创建最后一个分隔符之前的目录
func MkBaseDir(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
        _, err = os.Stat(filepath.Base(path))
        if os.IsNotExist(err) {
            os.MkdirAll(filepath.Dir(path), fs.ModeDir)
        }
	}
}

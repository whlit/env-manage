package util

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// 获取根目录 获取失败则直接退出程序
// 本方法以当前可执行文件所在的目录为bin目录为前提
// 注意使用
func GetRootDir() (string) {
	exePath, err := os.Executable()
	if err != nil {
        fmt.Println("获取根目录失败")
		os.Exit(1)
	}
    // 软件目录为 bin 根目录应该为上级目录
	return filepath.Dir(filepath.Dir(exePath))
}

// 获取当前可执行文件所在的目录
func GetExeDir() (string) {
	exePath, err := os.Executable()
	if err != nil {
        fmt.Println("获取可执行文件目录失败")
		os.Exit(1)
	}
	return filepath.Dir(exePath)
}

// 创建最后一个分隔符之前的目录
func MkBaseDir(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
        _, err = os.Stat(filepath.Dir(path))
        if os.IsNotExist(err) {
            os.MkdirAll(filepath.Dir(path), fs.ModeDir)
        }
	}
}

func GetConfigFilePath(name string) string {
    path := filepath.Join(GetRootDir(),"config", name)
    MkBaseDir(path)
	return path
}

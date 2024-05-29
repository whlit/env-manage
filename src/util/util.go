package util

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/whlit/env-manage/logger"
	"gopkg.in/yaml.v3"
)

// 获取根目录 获取失败则直接退出程序
// 本方法以当前可执行文件所在的目录为bin目录为前提
// 注意使用
func GetRootDir() string {
	exePath, err := os.Executable()
	if err != nil {
        logger.Error("获取根目录失败", err)
	}
	// 软件目录为 bin 根目录应该为上级目录
	return filepath.Dir(filepath.Dir(exePath))
}

// 获取当前可执行文件所在的目录
func GetExeDir() string {
	exePath, err := os.Executable()
	if err != nil {
        logger.Error("获取可执行文件目录失败", err)
	}
	return filepath.Dir(exePath)
}

// 获取当前可执行文件名称 不带后缀
func GetExeName() string {
	exePath, err := os.Executable()
	if err != nil {
        logger.Error("获取可执行文件目录失败", err)
	}
	name := filepath.Base(exePath)
	return strings.TrimSuffix(name, filepath.Ext(name))
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

// 获取配置文件路径
func GetConfigFilePath() string {
	path := filepath.Join(GetRootDir(), "config", strings.Join([]string{GetExeName(), ".yml"}, ""))
	MkBaseDir(path)
	return path
}

// 保存配置
func SaveConfig(config interface{}) {
	data, err := yaml.Marshal(config)
	if err != nil {
        logger.Warn("保存配置文件失败: ", err)
	}
	os.WriteFile(GetConfigFilePath(), data, 0644)
}

// 解压缩
func Unzip(zipPath, dir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dir, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dir, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dir)+string(os.PathSeparator)) {
            logger.Error("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

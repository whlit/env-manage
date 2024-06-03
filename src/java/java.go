package java

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/whlit/env-manage/core"
	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/util"
	"github.com/whlit/env-manage/version"
)

type JdkEnvManager struct {
	core.EnvManager
}

// 设置 JAVA_HOME
func (jdk *JdkEnvManager) Home(jhomePath string) (bool, error) {
	file, err := os.Stat(jhomePath)
	if err != nil {
		if strings.Contains(filepath.Base(jhomePath), ".") {
			return false, errors.New("JAVA_HOME需要是一个目录")
		}
		os.MkdirAll(filepath.Dir(jhomePath), fs.ModeDir)
		return true, nil
	}
	if !file.IsDir() {
		return false, errors.New("JAVA_HOME需要是一个目录")
	}
	dir, err := os.Open(jhomePath)
	if err != nil {
		return false, err
	}
	defer dir.Close()
	_, err = dir.Readdir(1)
	if err != io.EOF {
		return false, errors.New("目录必须为一个空目录,或者不存在的目录")
	}
	err = os.Remove(jhomePath)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (jdk *JdkEnvManager) Install() {
	dir := filepath.Join(util.GetRootDir(), "versions", jdk.Name)
	jdkVersion, err := version.Install(version.GetJdkVersions(), dir)
	if err != nil {
		logger.Error("安装失败", err)
		return
	}
	// 解压完成 开始配置
	logger.Info("正在添加到配置")
	dir = filepath.Join(dir, jdkVersion.GetVersionKey())
	dirs, err := os.ReadDir(dir)
	if err != nil {
		logger.Error("读取目录失败", err)
		return
	}
	if len(dirs) == 1 {
		dir = filepath.Join(dir, dirs[0].Name())
	}
	jdk.Add(jdkVersion.GetVersionKey(), dir)
	logger.Info("安装成功")
}

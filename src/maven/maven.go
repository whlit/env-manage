package maven

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

type MavenEnvManager struct {
	core.EnvManager
}

// 设置M2_HOME
func (maven *MavenEnvManager) Home(homePath string) (bool, error) {
	file, err := os.Stat(homePath)
	if err != nil {
		if strings.Contains(filepath.Base(homePath), ".") {
			return false, errors.New("M2_HOME需要是一个目录")
		}
		os.MkdirAll(filepath.Dir(homePath), fs.ModeDir)
		return true, nil
	}
	if !file.IsDir() {
		return false, errors.New("M2_HOME需要是一个目录")
	}
	dir, err := os.Open(homePath)
	if err != nil {
		return false, err
	}
	defer dir.Close()
	_, err = dir.Readdir(1)
	if err != io.EOF {
		return false, errors.New("目录必须为一个空目录,或者不存在的目录")
	}
	err = os.Remove(homePath)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (maven *MavenEnvManager) Install() {
	dir := filepath.Join(util.GetRootDir(), "versions", "maven")
	versions, err := version.Install(version.GetMavenVersions(), dir)
	if err != nil {
		logger.Error("安装失败", err)
	}
	// 解压完成 开始配置
	logger.Info("正在添加到配置")
	dir = filepath.Join(dir, versions.GetVersionKey())
	dirs, err := os.ReadDir(dir)
	if err != nil {
		logger.Error("读取目录失败", err)
	}
	if len(dirs) == 1 {
		dir = filepath.Join(dir, dirs[0].Name())
	}
	maven.Add(versions.GetVersionKey(), dir)
	logger.Info("安装成功")
}

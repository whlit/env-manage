package node

import (
	"os"
	"path/filepath"

	"github.com/whlit/env-manage/core"
	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/util"
	"github.com/whlit/env-manage/version"
)

type NodeEnvManager struct {
    core.EnvManager
}

func (node *NodeEnvManager) Install() {
    dir := filepath.Join(util.GetRootDir(), "versions", node.Name)
	jdkVersion, err := version.Install(version.GetNodeVersions(), dir)
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
	node.Add(jdkVersion.GetVersionKey(), dir)
	logger.Info("安装成功")
}



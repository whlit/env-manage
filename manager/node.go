package manager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/whlit/env-manage/core"
	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/util"
)

type NodeEnvManager struct {
    core.EnvManager
}

func NewManagerForNode() core.EnvManager {
	var m core.EnvManager
	m.Name = "node"
	m.Envs = make(map[string]map[string][]string)

	windowsEnv := make(map[string][]string)
	windowsEnv["NODE_HOME"] = []string{filepath.Join(util.GetRootDir(), core.GlobalConfig.RuntimeDir, m.Name)}
	windowsEnv["PATH"] = []string{"%NODE_HOME%"}
	m.Envs["windows"] = windowsEnv

	linuxEnv := make(map[string][]string)
	linuxEnv["NODE_HOME"] = []string{filepath.Join(util.GetRootDir(), core.GlobalConfig.RuntimeDir, m.Name)}
	linuxEnv["PATH"] = []string{"$NODE_HOME/bin"}
	m.Envs["linux"] = linuxEnv

	return m
}

func (m *NodeEnvManager) Install() {
	osMsg := m.getOsMsg()
    data, err := util.Get(fmt.Sprintf("https://raw.githubusercontent.com/whlit/versions/refs/heads/main/versions/node/node-%s-%s.version.json", osMsg.Os, osMsg.Arch))
    if err != nil {
        logger.Error("获取Node版本信息失败", err)
    }
	var versions map[string][]core.Version
	err = json.Unmarshal(data, &versions)
	if err != nil {
		logger.Error("解析Node版本信息失败", err, string(data))
	}
	version, ok := selectVersion(versions, osMsg.FileType)
	if !ok {
		return
	}
	version.App = m.Name

	err = version.Download()
	if err != nil {
		logger.Error("下载Node版本失败", err)
	}
    versionPath := version.GetVersionsPath()
	version.Path = filepath.Join(versionPath, version.FileName[:len(version.FileName)-len(version.FileType)-1])
	// 检查是否已经安装, 已安装则删除
	if util.FileExists(version.Path) {
		err = os.RemoveAll(versionPath)
		if err != nil {
			logger.Error("删除目录失败", err)
		}
	}
	if version.FileType == "zip" {
		err = util.Unzip(version.GetDownloadFilePath(), versionPath)
	} else if version.FileType == "tar.gz" {
		err = util.UnTarGz(version.GetDownloadFilePath(), versionPath)
	}
	if err != nil {
		logger.Error("解压失败：", err)
	}
	if mg, ok := core.GlobalConfig.Managers[m.Name]; ok {
		mg.Versions = append(mg.Versions, version)
		core.GlobalConfig.Managers[m.Name] = mg
		core.SaveConfig()
		logger.Info("安装成功")
	}
}


func (m *NodeEnvManager) getOsMsg() OsMsg {
	os := strings.ToLower(runtime.GOOS)
	switch os {
	case "windows", "win":
		os = "win"
	case "linux":
		os = "linux"
	case "darwin", "mac":
		os = "darwin"
	default:
		logger.Error("获取操作系统类型失败，或者不支持该系统类型", os)
	}

    arch := strings.ToLower(runtime.GOARCH)
    switch arch {
    case "amd64":
        arch = "x64"
    case "x86":
        arch = "x86"
    case "arm64":
        arch = "arm64"
    default:
        if os == "darwin" {
            arch = "any"
        } else {
            logger.Error("获取操作系统架构失败，或者不支持该系统架构", arch)
        }
    }

	fileType := "zip"
	if os == "windows" {
		fileType = "zip"
	} else if os == "linux" {
		fileType = "tar.gz"
	}
	return OsMsg{os, arch, fileType}
}

package manager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/whlit/env-manage/core"
	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/util"
)

type JdkEnvManager struct {
	core.EnvManager
}

type OsMsg struct {
	Os string
	Arch string
	FileType string
}

func NewManagerForJdk() core.EnvManager {
	var m core.EnvManager
	m.Name = "jdk"
	m.Envs = make(map[string]map[string][]string)

	windowsEnv := make(map[string][]string)
	windowsEnv["JAVA_HOME"] = []string{filepath.Join(util.GetRootDir(), core.GlobalConfig.RuntimeDir, m.Name)}
	windowsEnv["PATH"] = []string{"%JAVA_HOME%\\bin"}
	m.Envs["windows"] = windowsEnv

	linuxEnv := make(map[string][]string)
	linuxEnv["JAVA_HOME"] = []string{filepath.Join(util.GetRootDir(), core.GlobalConfig.RuntimeDir, m.Name)}
	linuxEnv["PATH"] = []string{"$JAVA_HOME/bin"}
	m.Envs["linux"] = linuxEnv

	return m
}

func (m *JdkEnvManager) Install() {
	osMsg := m.getOsMsg()
	data, err := util.Get(fmt.Sprintf("https://raw.githubusercontent.com/whlit/versions/refs/heads/main/versions/jdk/latest/jdk-%s-%s.version.json", osMsg.Os, osMsg.Arch))
	if err != nil {
		logger.Error("获取JDK版本信息失败", err)
	}
	var versions map[string][]core.Version
	err = json.Unmarshal(data, &versions)
	if err != nil {
		logger.Error("解析JDK版本信息失败", err, string(data))
	}
	version, ok := selectVersion(versions, osMsg.FileType)
	if !ok {
		return
	}
	version.App = m.Name

	err = version.Download()
	if err != nil {
		logger.Error("下载JDK版本失败", err)
	}
    versionPath := version.GetVersionsPath()
	version.Path = filepath.Join(versionPath, version.Version)
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
		mg.Versions = util.AppendOrUpdateSlice(mg.Versions, version, func(i int) bool { return mg.Versions[i].Version == version.Version })
		core.GlobalConfig.Managers[m.Name] = mg
		core.SaveConfig()
		logger.Info("安装成功")
	}
}

func selectVersion(versions map[string][]core.Version, fileType string) (core.Version, bool) {
	var items []core.Version
	for _, vs := range versions {
		for _, v := range vs {
			if v.FileType == fileType {
				items = append(items, v)
			}
		}
	}
	sort.SliceStable(items, func(i, j int) bool {
		return util.CompareVersion(items[i].Version, items[j].Version) > 0
	})
	return util.Select(func(v core.Version) string { return v.Version }, items...)
}

func (m *JdkEnvManager) getOsMsg() OsMsg {
	os := strings.ToLower(runtime.GOOS)
	switch os {
	case "windows", "win":
		os = "windows"
	case "linux":
		os = "linux"
	case "darwin", "mac":
		os = "mac"
	default:
		logger.Error("获取操作系统类型失败，或者不支持该系统类型", os)
	}

    arch := strings.ToLower(runtime.GOARCH)
    switch arch {
    case "amd64":
        arch = "x64"
    case "x86":
        arch = "x32"
    case "arm64":
        arch = "arm"
    case "aarch64":
        arch = "aarch64"
    default:
        logger.Error("获取操作系统架构失败，或者不支持该系统架构", arch)
    }

	fileType := "zip"
	if os == "windows" {
		fileType = "zip"
	} else if os == "linux" {
		fileType = "tar.gz"
	}
	return OsMsg{os, arch, fileType}
}

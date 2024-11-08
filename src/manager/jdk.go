package manager

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/whlit/env-manage/core"
	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/util"
)

type JdkEnvManager struct {
	core.EnvManager
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
    os, arch := m.getOsAndArch()
	data, err := util.Get(fmt.Sprintf("https://raw.githubusercontent.com/whlit/versions/refs/heads/main/versions/jdk/latest/jdk-%s-%s.version.json", os, arch))
	if err != nil {
		logger.Error("获取JDK版本信息失败", err)
	}
	var versions map[string][]core.Version
	err = json.Unmarshal(data, &versions)
	if err != nil {
		logger.Error("解析JDK版本信息失败", err, string(data))
	}
	version := m.selectVersion(versions)
	version.App = m.Name

	err = version.Download()
	if err != nil {
		logger.Error("下载JDK版本失败", err)
	}
    versionPath := version.GetVersionsPath()
	err = util.Unzip(version.GetDownloadFilePath(), versionPath)
	if err != nil {
		logger.Error("解压失败：", err)
	}
	if mg, ok := core.GlobalConfig.Managers[m.Name]; ok {
        version.Path = filepath.Join(versionPath, version.Version)
		mg.Versions = append(mg.Versions, version)
		core.GlobalConfig.Managers[m.Name] = mg
		core.SaveConfig()
		logger.Info("安装成功")
	}
}

func (m *JdkEnvManager) selectVersion(versions map[string][]core.Version) core.Version {
	var options []huh.Option[core.Version]
	for _, vs := range versions {
		for _, v := range vs {
			if v.FileType == "zip" {
				options = append(options, huh.NewOption(v.Version, v))
			}
		}
	}
	var version core.Version
	huh.NewSelect[core.Version]().Options(options...).Value(&version).Run()
	return version
}

func (m *JdkEnvManager) getOsAndArch() (string, string) {
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
    return os, arch
}
package manager

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/whlit/env-manage/core"
	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/util"
)

type MavenEnvManager struct {
    core.EnvManager
}

func init() {
    name := "maven"
    core.GlobalConfig.Register(name, func(em core.EnvManager) core.IEnvManager {return &MavenEnvManager{EnvManager: em}}, func() core.EnvManager {
        var m core.EnvManager
        m.Name = name
        m.Envs = make(map[string]map[string][]string)

        windowsEnv := make(map[string][]string)
        windowsEnv["M2_HOME"] = []string{filepath.Join(util.GetRootDir(), core.GlobalConfig.RuntimeDir, m.Name)}
        windowsEnv["PATH"] = []string{"%M2_HOME%\\bin"}
        m.Envs["windows"] = windowsEnv

        linuxEnv := make(map[string][]string)
        linuxEnv["M2_HOME"] = []string{filepath.Join(util.GetRootDir(), core.GlobalConfig.RuntimeDir, m.Name)}
        linuxEnv["PATH"] = []string{"$M2_HOME/bin"}
        m.Envs["linux"] = linuxEnv

        return m
    })
}

func (m *MavenEnvManager) Install() {
    data, err := util.Get("https://raw.githubusercontent.com/whlit/versions/refs/heads/main/versions/maven/maven.version.json")
    if err != nil {
        logger.Error("获取Maven版本信息失败", err)
    }
	var versions map[string][]core.Version
	err = json.Unmarshal(data, &versions)
	if err != nil {
		logger.Error("解析Maven版本信息失败", err, string(data))
	}
	version, ok := selectVersion(versions, "zip")
	if !ok {
		return
	}
	version.App = m.Name

	err = version.Download()
	if err != nil {
		logger.Error("下载Maven版本失败", err)
	}
    versionPath := version.GetVersionsPath()
	version.Path = filepath.Join(versionPath, version.FileName[:len(version.FileName)-len(version.FileType)-5])
    // 检查是否已经安装, 已安装则删除
	if util.FileExists(version.Path) {
		err = os.RemoveAll(versionPath)
		if err != nil {
			logger.Error("删除目录失败", err)
		}
	}
	err = util.Unzip(version.GetDownloadFilePath(), versionPath)
	if err != nil {
		logger.Error("解压失败：", err)
	}
    core.GlobalConfig.AddVersion(m.Name, version)
}

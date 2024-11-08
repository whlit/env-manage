package manager

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/charmbracelet/huh"
	"github.com/whlit/env-manage/core"
	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/util"
)

type MavenEnvManager struct {
    core.EnvManager
}

func NewManagerForMaven() core.EnvManager {
	var m core.EnvManager
	m.Name = "maven"
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
	version := m.selectVersion(versions)
	version.App = m.Name

	err = version.Download()
	if err != nil {
		logger.Error("下载Maven版本失败", err)
	}
    versionPath := version.GetVersionsPath()
	err = util.Unzip(version.GetDownloadFilePath(), versionPath)
	if err != nil {
		logger.Error("解压失败：", err)
	}
	if mg, ok := core.GlobalConfig.Managers[m.Name]; ok {
        version.Path = filepath.Join(versionPath, fmt.Sprintf("apache-maven-%s", version.Version))
		mg.Versions = append(mg.Versions, version)
		core.GlobalConfig.Managers[m.Name] = mg
		core.SaveConfig()
		logger.Info("安装成功")
	}
}
func (m *MavenEnvManager) selectVersion(versions map[string][]core.Version) core.Version {
	var options []huh.Option[core.Version]
	for _, vs := range versions {
		for _, v := range vs {
			if v.FileType == "zip" {
				options = append(options, huh.NewOption(v.Version, v))
			}
		}
	}
	var version core.Version
    sort.SliceStable(options, func(i, j int) bool {
        return util.CompareVersion(options[i].Value.Version, options[j].Value.Version) > 0
    })
	huh.NewSelect[core.Version]().Options(options...).Value(&version).Run()
	return version
}

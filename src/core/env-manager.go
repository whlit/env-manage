package core

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/util"
)

type IEnvManager interface {
	List()
	Add(name string, path string)
	Remove()
	Use()
	Install()
	CreateEnvs()
}

type EnvManager struct {
	Name     string                         `yaml:"name"`     // 名称 唯一 一般是软件名称
	Envs     map[string]map[string][]string `yaml:"envs"`     // 环境变量
	Versions []Version                      `yaml:"versions"` // 版本
}

// 列出已添加的版本
func (m *EnvManager) List() {
	table := util.NewTable("Version", "Path")
	for _, v := range m.Versions {
		table.Add(map[string]string{
			"Version": v.Version,
			"Path":    v.Path,
		})
	}
	table.Printf()
}

// 添加版本
func (m *EnvManager) Add(name string, path string) {
	file, err := os.Stat(path)
	if err != nil {
		logger.Error("目录不存在: ", path)
	}
	if !file.IsDir() {
		logger.Error("路径不是目录: ", path)
	}
	for _, v := range m.Versions {
		if v.Version == name {
			logger.Warn("版本已存在: ", name, " -> ", v.Path)
			return
		}
	}
	version := &Version{}
	version.Version = name
	version.Path = path
	m.Versions = append(m.Versions, *version)
	SaveConfig()
}

// 移除版本

func (m *EnvManager) Remove() {
	if m.Versions == nil || len(m.Versions) == 0 {
		logger.Info("未添加任何版本")
		return
	}
	// 选择版本
	var version Version
	var options []huh.Option[Version]
	for _, v := range m.Versions {
		options = append(options, huh.NewOption(v.Version, v))
	}
	huh.NewSelect[Version]().Options(options...).Value(&version).Run()
	var confirm bool
	huh.NewConfirm().Title(strings.Join([]string{"确认删除 ", version.Version, " ?"}, "")).Value(&confirm).Run()
	if confirm {
		vs := m.Versions[:0]
		for _, v := range m.Versions {
			if v.Version != version.Version {
				vs = append(vs, v)
			}
		}
		m.Versions = vs
		SaveConfig()
	}
}

// 使用版本
func (m *EnvManager) Use() {
	if m.Versions == nil || len(m.Versions) == 0 {
		logger.Info("未添加任何版本")
		return
	}
	// 选择版本
	var version Version
	var options []huh.Option[Version]
	for _, v := range m.Versions {
		options = append(options, huh.NewOption(v.Version, v))
	}
	huh.NewSelect[Version]().Options(options...).Value(&version).Run()

	path := filepath.Join(util.GetRootDir(), GlobalConfig.RuntimeDir, m.Name)

	if util.FileExists(path) {
		os.Remove(path)
	}
	util.CreateLink(path, version.Path)
	logger.Info("成功切换版本为", version.Version)
}

// 安装
func (m *EnvManager) Install() {

}

// 创建环境变量
func (m *EnvManager) CreateEnvs() {

}

func RegisterEnvManager(manager EnvManager) {
	_, ok := GlobalConfig.Managers[manager.Name]
	if ok {
		return
	}
	GlobalConfig.Managers[manager.Name] = manager
	SaveConfig()
}

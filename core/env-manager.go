package core

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/util"
)

type IEnvManager interface {
	List()
	Add(name string, path string)
	Remove()
	Use()
	Install()
	InitEnvs()
}

type EnvManager struct {
	Name     string                         `yaml:"name"`     // 名称 唯一 一般是软件名称
	Envs     map[string]map[string][]string `yaml:"envs"`     // 环境变量
	Versions []Version                      `yaml:"versions"` // 版本
    Used     string                         `yaml:"used"`     // 当前使用的版本
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
			logger.Error("版本已存在: ", name, " -> ", v.Path)
			return
		}
	}
	version := &Version{}
	version.Version = name
	version.Path = path
    GlobalConfig.AddVersion(m.Name, *version)
}

// 移除版本

func (m *EnvManager) Remove() {
	if len(m.Versions) == 0 {
		logger.Info("未添加任何版本")
		return
	}
	// 选择版本
	if version, ok := util.SelectWithConfirm(func (v Version) string { return v.Version }, m.Versions...); ok {
        GlobalConfig.RemoveVersion(m.Name, version)
	}
}

// 使用版本
func (m *EnvManager) Use() {
	if len(m.Versions) == 0 {
		logger.Info("未添加任何版本")
		return
	}
	// 选择版本
	version, ok := util.Select(func(v Version) string { return v.Version }, m.Versions...)
	if !ok {
		return
	}
	path := filepath.Join(util.GetRootDir(), GlobalConfig.RuntimeDir, m.Name)

	if util.FileExists(path) {
		os.Remove(path)
	}
	err := util.CreateLink(path, version.Path)
    GlobalConfig.SetUsed(m.Name, version)
	if err != nil {
		logger.Error("创建链接失败：", err)
	}
}

// 安装
func (m *EnvManager) Install() {
    fmt.Println("不支持在线安装: ", m.Name)
}

// 创建环境变量
func (m *EnvManager) InitEnvs() {
    if _, ok := m.Envs[runtime.GOOS]; ok {
		util.SetEnvs(m.Envs[runtime.GOOS])
		return
    }
    logger.Error("暂不支持自动创建该系统环境变量，请手动设置")
}

package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/util"
)

type IEnvManager interface {
	List()
	Add(name string, path string)
	Remove()
	Use() (string, string, error)
	Install()
	GetEnvs() map[string][]string
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
    if mg, ok := GlobalConfig.Managers[m.Name]; ok {
        mg.Versions = append(mg.Versions, *version)
        GlobalConfig.Managers[m.Name] = mg
        SaveConfig()
    }
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
        if mg, ok := GlobalConfig.Managers[m.Name]; ok {
            var vs []Version
            for _, v := range mg.Versions {
                if v.Version != version.Version {
                    vs = append(vs, v)
                }
            }
            mg.Versions = vs
            GlobalConfig.Managers[m.Name] = mg
            SaveConfig()
        }
	}
}

// 使用版本
func (m *EnvManager) Use() (string, string, error) {
	if m.Versions == nil || len(m.Versions) == 0 {
		return "", "", errors.New("未添加任何版本")
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
    return path, version.Path, nil
}

// 安装
func (m *EnvManager) Install() {
    fmt.Println("不支持在线安装: ", m.Name)
}

// 创建环境变量
func (m *EnvManager) GetEnvs() map[string][]string {
    if _, ok := m.Envs[runtime.GOOS]; ok {
        return m.Envs[runtime.GOOS]
    }
    logger.Error("暂不支持自动创建该系统环境变量，请手动设置")
    return nil
}

func RegisterEnvManager(manager EnvManager) {
	_, ok := GlobalConfig.Managers[manager.Name]
	if ok {
		return
	}
	GlobalConfig.Managers[manager.Name] = manager
	SaveConfig()
}

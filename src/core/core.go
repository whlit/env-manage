package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/whlit/env-manage/cmd"
	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/util"
	"golang.org/x/exp/maps"
)

type IEnvManager interface {
	List()
	Add(name string, path string)
	Remove()
	Use()
}

type IInstallable interface {
	Install()
}

type EnvManager struct {
	Name         string            `yaml:"name"`           // 名称 唯一 一般是软件名称
	EnvName      string            `yaml:"env_name"`       // 环境变量名称
	EnvValue     string            `yaml:"env_value"`      // 环境变量值
	EnvPathValue string            `yaml:"env_path_value"` // 添加到Path的环境变量值
	Versions     map[string]string `yaml:"versions"`       // 版本列表
}

// 列出已添加的版本
func (m *EnvManager) List() {
	var used string
	if m.EnvValue != "" {
		used, _ = os.Readlink(m.EnvValue)
	}
	table := util.Table{
		Columns: []string{"Version", "Path"},
		Selected: func(row map[string]string) bool {
			return row["Path"] == used
		},
	}
	for k, v := range m.Versions {
		table.Add(map[string]string{
			"Version": k,
			"Path":    v,
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
	m.Versions[name] = path
}

// 移除版本

func (m *EnvManager) Remove() {
	if m.Versions == nil || len(m.Versions) == 0 {
		logger.Info("未添加任何版本")
		return
	}
	// 选择版本
	var name string
	huh.NewSelect[string]().Options(huh.NewOptions(maps.Keys(m.Versions)...)...).Value(&name).Run()
	if m.Versions[name] == "" {
		logger.Info("版本不存在")
		return
	}
	var confirm bool
	huh.NewConfirm().Title(strings.Join([]string{"确认删除 ", name, " ?"}, "")).Value(&confirm).Run()
	if confirm {
		delete(m.Versions, name)
	}
}

// 使用版本
func (m *EnvManager) Use() {
	if m.Versions == nil || len(m.Versions) == 0 {
		logger.Info("未添加任何版本")
		return
	}
	// 选择版本
	var name string
	huh.NewSelect[string]().Options(huh.NewOptions(maps.Keys(m.Versions)...)...).Value(&name).Run()
	if m.Versions[name] == "" {
		logger.Info("版本不存在")
		return
	}
	// 设置环境变量
	if m.EnvValue == "" {
		m.EnvValue = filepath.Join(util.GetRootDir(), "runtime", m.Name)
	}
	cmd.SetEnvironmentValue(m.EnvName, m.EnvValue)
	// 添加到PATH
	if m.EnvPathValue == "" {
		m.EnvPathValue = fmt.Sprintf("%%%s%%\\bin", m.EnvName)
	}
	cmd.AddToPath(m.EnvPathValue)
	cmd.CreateLink(m.EnvValue, m.Versions[name])
	logger.Info("成功切换版本为", name)
}

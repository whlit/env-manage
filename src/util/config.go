package util

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/whlit/env-manage/logger"
	"gopkg.in/yaml.v3"
)

type Config[T any] struct {
	Data T
	name string
}

// 加载配置
func (c *Config[T]) Load() {
	if FileExists(c.getFilePath()) {
		file, err := os.ReadFile(c.getFilePath())
		if err != nil {
			logger.Error("配置文件读取失败,", c.getFilePath(), err)
		}
		yaml.Unmarshal(file, &c.Data)
	} else {
		os.Create(c.getFilePath())
	}
}

// 保存配置
func (c *Config[T]) Save() {
	data, err := yaml.Marshal(c.Data)
	if err != nil {
		logger.Warn("保存配置文件失败: ", err)
	}
	os.WriteFile(c.getFilePath(), data, 0644)
}

// 新建一个配置
func NewConfig[T any](name string, config T) *Config[T] {
	return &Config[T]{
		Data: config,
		name: name,
	}
}

// 获取配置文件路径
func (c *Config[T]) getFilePath() string {
	path := filepath.Join(GetRootDir(), "config", strings.Join([]string{c.name, ".yml"}, ""))
	MkBaseDir(path)
	return path
}

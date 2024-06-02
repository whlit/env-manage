package core

import (
	"os"
	"path/filepath"

	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/util"
	"gopkg.in/yaml.v3"
)

type Config[T any] struct {
	Data T
	name string
}

// 加载配置
func (c *Config[T]) Load() {
	file, err := os.ReadFile(getFilePath())
	if err != nil {
		logger.Error("配置文件读取失败,", getFilePath(), err)
	}
	var config map[string]T
	yaml.Unmarshal(file, &config)
	if config == nil {
		return
	}
	if _, ok := config[c.name]; ok {
		c.Data = config[c.name]
	}
}

// 保存配置
func (c *Config[T]) Save() {
	var config map[string]any
	file, err := os.ReadFile(getFilePath())
	if err != nil {
		logger.Error("配置文件读取失败,", getFilePath(), err)
	}
	yaml.Unmarshal(file, &config)
	if config == nil {
		config = make(map[string]any)
	}
	config[c.name] = c.Data
	data, err := yaml.Marshal(config)
	if err != nil {
		logger.Warn("保存配置文件失败: ", err)
	}
	os.WriteFile(getFilePath(), data, 0644)
}

// 新建一个配置
func NewConfig[T any](name string, config T) *Config[T] {
	return &Config[T]{
		Data: config,
		name: name,
	}
}

// 获取配置文件路径
func getFilePath() string {
	path := filepath.Join(util.GetRootDir(), "config", "config.yml")
	if !util.FileExists(path) {
		util.MkBaseDir(path)
		os.Create(path)
	}
	return path
}

// 获取所有配置名称
func GetConfigs() map[string]any {
	var config map[string]any
	file, err := os.ReadFile(getFilePath())
	if err != nil {
		logger.Error("配置文件读取失败,", getFilePath(), err)
	}
	yaml.Unmarshal(file, &config)
	if config == nil {
		config = make(map[string]any)
	}
	return config
}

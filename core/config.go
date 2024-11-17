package core

import (
	"os"
	"path/filepath"

	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/util"
	"gopkg.in/yaml.v3"
)

const (
	CONFIG_FILE_NAME     = "config.yml"
	CONFIG_DIR           = "config"
	DEFAULT_DOWNLOAD_DIR = "download"
	DEFAULT_VERSIONS_DIR = "versions"
	DEFAULT_RUNTIME_DIR  = "runtime"
)

type Config struct {
	VersionsDir string       `yaml:"versions_dir"`
	RuntimeDir  string       `yaml:"runtime_dir"`
	DownloadDir string       `yaml:"download_dir"`
	LoggerLevel string       `yaml:"logger_level"`
	Managers    []EnvManager `yaml:"managers"`
}

var GlobalConfig = Config{}
var managerMap = make(map[string]IEnvManager)

func init() {
	root := util.GetRootDir()
	path := filepath.Join(root, CONFIG_DIR, CONFIG_FILE_NAME)
	// 读取配置
	if util.FileExists(path) {
		file, err := os.ReadFile(path)
		if err != nil {
			logger.Error("配置文件读取失败,", path, err)
			os.Exit(1)
		}
		err = yaml.Unmarshal(file, &GlobalConfig)
		if err != nil {
			logger.Error("配置文件解析失败,", path, err)
			os.Exit(1)
		}
		for _, manager := range GlobalConfig.Managers {
			managerMap[manager.Name] = &manager
		}
		if GlobalConfig.LoggerLevel != "" {
			logger.SetLevel(GlobalConfig.LoggerLevel)
		}
		return
	}
	// 初始化配置 并创建配置文件
	util.MkBaseDir(path)
	os.Create(path)
	GlobalConfig = Config{
		VersionsDir: DEFAULT_VERSIONS_DIR,
		RuntimeDir:  DEFAULT_RUNTIME_DIR,
		DownloadDir: DEFAULT_DOWNLOAD_DIR,
		LoggerLevel: logger.INFO.String(),
		Managers:    make([]EnvManager, 0),
	}
	// 创建默认文件夹
	os.MkdirAll(filepath.Join(root, GlobalConfig.VersionsDir), 00755)
	os.MkdirAll(filepath.Join(root, GlobalConfig.RuntimeDir), 00755)
	os.MkdirAll(filepath.Join(root, GlobalConfig.DownloadDir), 00755)
	GlobalConfig.Save()
}

func (c *Config) Save() {
	data, err := yaml.Marshal(c)
	if err != nil {
		logger.Error("序列化配置失败: ", err)
		return
	}
	err = os.WriteFile(filepath.Join(util.GetRootDir(), CONFIG_DIR, CONFIG_FILE_NAME), data, 0644)
	if err != nil {
		logger.Error("保存配置文件失败: ", err)
		return
	}
}

func (c *Config) AddVersion(name string, version Version) {
	for i, m := range c.Managers {
		if m.Name == name {
			for j, v := range m.Versions {
				if v.Version == version.Version {
					c.Managers[i].Versions[j] = version
					c.Save()
					return
				}
			}
			c.Managers[i].Versions = append(c.Managers[i].Versions, version)
			c.Save()
			return
		}
	}
}

func (c *Config) RemoveVersion(name string, version Version) {
	for i, m := range c.Managers {
		if m.Name == name {
			for j, v := range m.Versions {
				if v.Version == version.Version {
					c.Managers[i].Versions = append(c.Managers[i].Versions[:j], c.Managers[i].Versions[j+1:]...)
					c.Save()
					return
				}
			}
		}
	}
}

func (c *Config) SetUsed(name string, version Version) {
	for i, m := range c.Managers {
		if m.Name == name {
			for j, v := range m.Versions {
				if v.Version == version.Version {
					c.Managers[i].Used = c.Managers[i].Versions[j].Version
					c.Save()
					return
				}
			}
		}
	}
}

func (c *Config) Register(name string, toManager func(EnvManager) IEnvManager, getDefault func() EnvManager) {
	if _, ok := managerMap[name]; ok {
		for _, manager := range c.Managers {
			if manager.Name == name {
				managerMap[name] = toManager(manager)
				return
			}
		}
	}
	defaultManager := getDefault()
	c.Managers = append(c.Managers, defaultManager)
	managerMap[name] = toManager(defaultManager)
	c.Save()
}

func (c *Config) GetManager(name string) IEnvManager {
	if manager, ok := managerMap[name]; ok {
		return manager
	}
	return nil
}

// 保存配置
func SaveConfig() error {
	data, err := yaml.Marshal(GlobalConfig)
	if err != nil {
		logger.Warn("序列化配置失败: ", err)
		return err
	}
	err = os.WriteFile(filepath.Join(util.GetRootDir(), CONFIG_DIR, CONFIG_FILE_NAME), data, 0644)
	if err != nil {
		logger.Warn("保存配置文件失败: ", err)
	}
	return err
}

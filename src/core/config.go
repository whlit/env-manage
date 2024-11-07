package core

import (
	"io/fs"
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
	VersionsDir string                `yaml:"versions_dir"`
	RuntimeDir  string                `yaml:"runtime_dir"`
	DownloadDir string                `yaml:"download_dir"`
	Managers    map[string]EnvManager `yaml:"managers"`
}

var GlobalConfig = Config{}

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
		return
	}
	// 初始化配置 并创建配置文件
	util.MkBaseDir(path)
	os.Create(path)
	GlobalConfig = Config{
		VersionsDir: DEFAULT_VERSIONS_DIR,
		RuntimeDir:  DEFAULT_RUNTIME_DIR,
		DownloadDir: DEFAULT_DOWNLOAD_DIR,
		Managers:    map[string]EnvManager{},
	}
	// 创建默认文件夹
	os.MkdirAll(filepath.Join(root, GlobalConfig.VersionsDir), fs.ModeDir)
	os.MkdirAll(filepath.Join(root, GlobalConfig.RuntimeDir), fs.ModeDir)
	os.MkdirAll(filepath.Join(root, GlobalConfig.DownloadDir), fs.ModeDir)
	SaveConfig()
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

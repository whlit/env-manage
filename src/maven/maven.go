package maven

import (
	"path/filepath"

	"github.com/whlit/env-manage/cmd"
	"github.com/whlit/env-manage/util"
)

type MavenConfig struct {
	Mavens map[string]string `yaml:"Mavens"`
	M2Home string            `yaml:"M2Home"`
}

var config = util.NewConfig("maven", &MavenConfig{
	Mavens: make(map[string]string),
	M2Home: filepath.Join(util.GetRootDir(), "runtime", "maven"),
})

// 初始化
func InitConfig() {
	// 加载配置文件
	config.Load()
	// 写入M2_HOME环境变量
	cmd.SetEnvironmentValue("M2_HOME", config.Data.M2Home)
	// 初始化JAVA_HOME到Path
	cmd.AddToPath("%M2_HOME%\\bin")
}

func List() {

}

func Add(version string, mavenPath string) {

}

func Remove(name string) {

}

func Use() {

}

func Home(homePath string) {

}

func Install() {

}

func Help() {

}

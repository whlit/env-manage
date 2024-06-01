package maven

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/whlit/env-manage/core"
	"github.com/whlit/env-manage/util"
	"github.com/whlit/env-manage/version"
)

var vm core.EnvManager = core.EnvManager{
	EnvName:  "JAVA_HOME",
	Versions: make(map[string]string),
	Name:     "jdk",
}

var config = core.NewConfig("maven", &vm)

// 初始化
func InitConfig() {
	// 加载配置文件
	config.Load()
}

// 列出所有已安装Maven
func List() {
	vm.List()
}

// 添加Maven
func Add(version string, mavenPath string) {
	vm.Add(version, mavenPath)
	config.Save()
}

// 移除Maven
func Remove() {
	vm.Remove()
	config.Save()
}

// 切换Maven
func Use() {
	vm.Use()
}

// 设置M2_HOME
func Home(homePath string) {
	if vm.EnvValue == homePath {
		return
	}
	defer config.Save()
	file, err := os.Stat(homePath)
	if err != nil {
		if strings.Contains(filepath.Base(homePath), ".") {
			fmt.Println("M2_HOME需要是一个目录")
			return
		}
		os.MkdirAll(filepath.Dir(homePath), fs.ModeDir)
		vm.EnvValue = homePath
		return
	}
	if !file.IsDir() {
		fmt.Println("M2_HOME需要是一个目录")
		return
	}
	dir, err := os.Open(homePath)
	if err != nil {
		fmt.Println("获取目录失败")
		return
	}
	defer dir.Close()
	_, err = dir.Readdir(1)
	if err != io.EOF {
		fmt.Println("目录必须为一个空目录,或者不存在的目录")
		return
	}
	err = os.Remove(homePath)
	if err != nil {
		fmt.Println("删除目录失败", err)
	}
	vm.EnvValue = homePath
}

func Install() {
	dir := filepath.Join(util.GetRootDir(), "versions", "maven")
	maven, err := version.Install(version.GetMavenVersions(), dir)
	if err != nil {
		fmt.Println("安装失败", err)
		return
	}
	// 解压完成 开始配置
	fmt.Println("正在添加到配置")
	dir = filepath.Join(dir, maven.GetVersionKey())
	dirs, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("读取目录失败", err)
		return
	}
	if len(dirs) == 1 {
		dir = filepath.Join(dir, dirs[0].Name())
	}
	Add(maven.GetVersionKey(), dir)
	fmt.Println("安装成功")
}

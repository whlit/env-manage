package java

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

var config = core.NewConfig("jdk", &core.EnvManager{
	EnvName:  "JAVA_HOME",
	Versions: make(map[string]string),
	Name:     "jdk",
})

type JdkEnvManager struct {
	core.EnvManager
}

// 初始化
func InitConfig() {
	// 加载配置文件
	config.Load()
}

// 列出所有已安装 JDK
func List() {
	config.Data.List()
}

// 添加 JDK
func Add(version string, jpath string) {
	config.Data.Add(version, jpath)
	config.Save()
}

// 移除 JDK
func Remove() {
	config.Data.Remove()
	config.Save()
}

// 切换 JDK
func Use() {
	config.Data.Use()
	config.Save()
}

// 设置 JAVA_HOME
func Home(jhomePath string) {
	if config.Data.EnvValue == jhomePath {
		return
	}
	defer config.Save()
	file, err := os.Stat(jhomePath)
	if err != nil {
		if strings.Contains(filepath.Base(jhomePath), ".") {
			fmt.Println("JAVA_HOME需要是一个目录")
			return
		}
		os.MkdirAll(filepath.Dir(jhomePath), fs.ModeDir)
		config.Data.EnvValue = jhomePath
		return
	}
	if !file.IsDir() {
		fmt.Println("JAVA_HOME需要是一个目录")
		return
	}
	dir, err := os.Open(jhomePath)
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
	err = os.Remove(jhomePath)
	if err != nil {
		fmt.Println("删除目录失败", err)
	}
	config.Data.EnvValue = jhomePath
}

func Install() {
	dir := filepath.Join(util.GetRootDir(), "versions", "jdk")
	jdk, err := version.Install(version.GetJdkVersions(), dir)
	if err != nil {
		fmt.Println("安装失败", err)
		return
	}
	// 解压完成 开始配置
	fmt.Println("正在添加到配置")
	dir = filepath.Join(dir, jdk.GetVersionKey())
	dirs, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("读取目录失败", err)
		return
	}
	if len(dirs) == 1 {
		dir = filepath.Join(dir, dirs[0].Name())
	}
	Add(jdk.GetVersionKey(), dir)
	fmt.Println("安装成功")
}

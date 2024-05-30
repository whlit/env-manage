package maven

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/whlit/env-manage/cmd"
	"github.com/whlit/env-manage/util"
	"github.com/whlit/env-manage/version"
	"golang.org/x/exp/maps"
	"golang.org/x/text/encoding/simplifiedchinese"
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
	// 初始化M2_HOME到Path
	cmd.AddToPath("%M2_HOME%\\bin")
}

// 列出所有已安装Maven
func List() {
    var used string
	if config.Data.M2Home != "" {
		used, _ = os.Readlink(config.Data.M2Home)
	}
	table := util.Table{
		Columns: []string{"Version", "Path"},
		Selected: func(row map[string]string) bool {
			return row["Path"] == used
		},
	}
	for k, v := range config.Data.Mavens {
		table.Add(map[string]string{
			"Version": k,
			"Path":    v,
		})
	}
	table.Printf()
}

// 添加Maven
func Add(version string, mavenPath string) {
    if !util.FileExists(mavenPath) {
		fmt.Println("路径不存在")
		return
	}
	if !util.FileExists(filepath.Join(mavenPath, "bin/mvn.cmd")) {
		fmt.Println("路径不是 Maven 路径", mavenPath)
	}
	config.Data.Mavens[version] = mavenPath
	config.Save()
}

// 移除Maven
func Remove(name string) {
    delete(config.Data.Mavens, name)
	config.Save()
}

// 切换Maven
func Use() {
    if config.Data.M2Home == "" {
		fmt.Println("请先设置 M2_HOME. 使用命令 home <path>")
		return
	}
	if config.Data.Mavens == nil || len(config.Data.Mavens) == 0 {
		fmt.Println("未添加任何Maven版本")
		return
	}
	var name string
	huh.NewSelect[string]().Options(huh.NewOptions(maps.Keys(config.Data.Mavens)...)...).Value(&name).Run()
	if config.Data.Mavens[name] == "" {
		fmt.Println("Maven 版本不存在")
		return
	}
	home, _ := os.Lstat(config.Data.M2Home)
	if home != nil {
		cmd.ElevatedRun("rmdir", filepath.Clean(config.Data.M2Home))
	}
    util.MkBaseDir(filepath.Clean(config.Data.M2Home))
	_, err := cmd.ElevatedRun("mklink", "/D", filepath.Clean(config.Data.M2Home), filepath.Clean(config.Data.Mavens[name]))
	if err != nil {
		errr, _ := simplifiedchinese.GB18030.NewDecoder().String(err.Error())
		fmt.Println(errr)
		return
	}
	fmt.Println("成功切换Maven版本为", name)
}

// 设置M2_HOME
func Home(homePath string) {
    if config.Data.M2Home == homePath {
		return
	}
	file, err := os.Stat(homePath)
	if err != nil {
		if strings.Contains(filepath.Base(homePath), ".") {
			fmt.Println("M2_HOME需要是一个目录")
			return
		}
		os.MkdirAll(filepath.Dir(homePath), fs.ModeDir)
		saveM2Home(homePath)
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
	saveM2Home(homePath)
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

func saveM2Home(homePath string) {
	config.Data.M2Home = homePath
	config.Save()
	cmd.SetEnvironmentValue("M2_HOME", homePath)
	fmt.Println("设置M2_HOME成功,需要重启终端生效")
}

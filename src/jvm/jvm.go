package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/text/encoding/simplifiedchinese"
	"gopkg.in/yaml.v3"

	"github.com/charmbracelet/huh"
	"github.com/whlit/env-manage/cmd"
	"github.com/whlit/env-manage/util"
)

type Config struct {
	Jdks  map[string]string `yaml:"Jdks"`
	Jhome string            `yaml:"Jhome"`
	Root  string            `yaml:"Root"`
}

var config = &Config{
	Jhome: "",
	Jdks:  make(map[string]string),
	Root:  "",
}

func main() {
	args := os.Args
	action := ""
	if len(args) > 1 {
		action = args[1]
	}

	switch action {
	case "list":
		list()
	case "add":
		add(args[2], args[3])
	case "rm":
		remove(args[2])
	case "use":
		use()
	case "home":
		home(args[2])
	case "install":
		install()
	default:
		help()
	}
}

// 列出所有已安装 JDK
func list() {
	var used string
	if config.Jhome != "" {
		used, _ = os.Readlink(config.Jhome)
	}
	table := util.Table{
		Columns: []string{"Version", "Path"},
		Selected: func(row map[string]string) bool {
			return row["Path"] == used
		},
	}
	for k, v := range config.Jdks {
		table.Add(map[string]string{
			"Version": k,
			"Path":    v,
		})
	}
	table.Print()

}

// 添加 JDK
func add(version string, jpath string) {
	if !fileExists(jpath) {
		fmt.Println("路径不存在")
		return
	}
	if !fileExists(path.Join(jpath, "bin/java.exe")) {
		fmt.Println("路径不是 JDK 路径")
	}
	config.Jdks[version] = jpath
	util.SaveConfig(config)
}

// 移除 JDK
func remove(name string) {
	delete(config.Jdks, name)
	util.SaveConfig(config)
}

// 切换 JDK
func use() {
	if config.Jhome == "" {
		fmt.Println("请先设置 JAVA_HOME. 使用命令 jvm home <path>")
		return
	}
	if config.Jdks == nil || len(config.Jdks) == 0 {
		fmt.Println("未添加任何JDK版本")
		return
	}
	var name string
	huh.NewSelect[string]().Options(huh.NewOptions(maps.Keys(config.Jdks)...)...).Value(&name).Run()
	if config.Jdks[name] == "" {
		fmt.Println("JDK 版本不存在")
		return
	}
	home, _ := os.Lstat(config.Jhome)
	if home != nil {
		cmd.ElevatedRun("rmdir", filepath.Clean(config.Jhome))
	}
	_, err := cmd.ElevatedRun("mklink", "/D", filepath.Clean(config.Jhome), filepath.Clean(config.Jdks[name]))
	if err != nil {
		errr, _ := simplifiedchinese.GB18030.NewDecoder().String(err.Error())
		fmt.Println(errr)
		return
	}
	fmt.Println("成功切换JAVA版本为", name)
}

// 设置 JAVA_HOME
func home(jhomePath string) {
	if config.Jhome == jhomePath {
		return
	}
	file, err := os.Stat(jhomePath)
	if err != nil {
		if strings.Contains(filepath.Base(jhomePath), ".") {
			fmt.Println("JAVA_HOME需要是一个目录")
			return
		}
		os.MkdirAll(filepath.Dir(jhomePath), fs.ModeDir)
		setJavaHome(jhomePath)
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
	setJavaHome(jhomePath)
}

func install() {
	var version string
	huh.NewSelect[string]().Options(huh.NewOptions(getDownloadVersions()...)...).Value(&version).Run()
	fmt.Println("开始安装JDK", version)
}

func setJavaHome(jhome string) {
	config.Jhome = jhome
	cmd.SetEnvironmentValue("JAVA_HOME", jhome)
	util.SaveConfig(config)
	fmt.Println("设置JAVA_HOME成功,需要重启终端生效")
}

// 初始化
func init() {
	// 加载配置文件
	loadConfig()
	// 初始化JAVA_HOME到Path
	cmd.AddToPath("%JAVA_HOME%\\bin")
}

func loadConfig() {
	root := util.GetRootDir()
	var configFile = util.GetConfigFilePath()
	// 读取配置文件
	if fileExists(configFile) {
		file, err := os.ReadFile(configFile)
		if err != nil {
			fmt.Println("读取配置文件失败")
		}
		var yamlData = &Config{}
		yaml.Unmarshal(file, &yamlData)
		// 设置 JDK 列表
		if yamlData.Jdks != nil {
			config.Jdks = yamlData.Jdks
		}
		// 设置 JAVA_HOME
		config.Jhome = yamlData.Jhome
		// 设置根目录
		if yamlData.Root != "" {
			config.Root = yamlData.Root
		} else {
			config.Root = root
		}
		return
	} else {
		os.Create(configFile)
		config.Root = root
		config.Jhome = path.Join(root, "runtime/jdk")
		util.MkBaseDir(config.Jhome)
		util.SaveConfig(config)
	}
}

func help() {
	fmt.Println("jvm add <name> <path>           Add a JDK")
	fmt.Println("jvm rm <name>                   Remove a JDK")
	fmt.Println("jvm list                        List all installed JDKs")
	fmt.Println("jvm use                         Select And Use a JDK")
	fmt.Println("jvm home <path>                 Set the path of JAVA_HOME")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

type DownloadInfo struct {
	Version      string
	Url          string
	CheckCodeUrl string
	CheckType    string
}

func getDownloads() []DownloadInfo {
	return []DownloadInfo{
		{
			Version:      "jdk22",
			Url:          "https://download.oracle.com/java/22/latest/jdk-22_windows-x64_bin.zip",
			CheckType:    "sha256",
			CheckCodeUrl: "https://download.oracle.com/java/22/latest/jdk-22_windows-x64_bin.zip.sha256",
		},
		{
			Version:      "jdk21",
			Url:          "https://download.oracle.com/java/21/latest/jdk-21_windows-x64_bin.zip",
			CheckType:    "sha256",
			CheckCodeUrl: "https://download.oracle.com/java/21/latest/jdk-21_windows-x64_bin.zip.sha256",
		},
		{
			Version:      "jdk17",
			Url:          "https://download.oracle.com/java/17/latest/jdk-17_windows-x64_bin.zip",
			CheckType:    "sha256",
			CheckCodeUrl: "https://download.oracle.com/java/17/latest/jdk-17_windows-x64_bin.zip.sha256",
		},
	}
}

func getDownloadVersions() []string {
	var versions []string
	for _, d := range getDownloads() {
		versions = append(versions, d.Version)
	}
	return versions
}

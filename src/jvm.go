package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"gopkg.in/yaml.v3"

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
		use(args[2])
	case "home":
		home(args[2])
    case "global":
        global(args[2])
	default:
		help()
	}
}

func list() {
	var used string
	if config.Jhome != "" {
		used, _ = os.Readlink(config.Jhome)
	}
	var ks, vs int
	for k, v := range config.Jdks {
		ks = max(ks, len(k))
		vs = max(vs, len(v))
	}
	kf := " %s %-" + strconv.Itoa(ks) + "s   %-" + strconv.Itoa(vs) + "s\n"
	for k, v := range config.Jdks {
		if used == v {
			fmt.Printf(kf, "*", k, v)
		} else {
			fmt.Printf(kf, " ", k, v)
		}
	}
}

func add(version string, jpath string) {
	if !fileExists(jpath) {
		fmt.Println("路径不存在")
		return
	}
	if !fileExists(path.Join(jpath, "bin/java.exe")) {
		fmt.Println("路径不是 JDK 路径")
	}
	config.Jdks[version] = jpath
	writeConfig()
}

func remove(name string) {
	delete(config.Jdks, name)
	writeConfig()
}

func use(name string) {
	if config.Jhome == "" {
		fmt.Println("请先设置 JAVA_HOME. 使用命令 jvm home <path>")
		return
	}
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
	}
}

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

func global(action string) {
	switch action {
	case "install":
        install()
	case "uninstall":
        unInstall()
	default:
		help()
	}
}

func install() {
    root, err := util.GetRootDir()
    if err != nil {
		fmt.Println("获取根目录失败")
		os.Exit(1)
	}
    cmd.AddToPath(root)
    fmt.Println("安装成功,请重新打开终端使用")
}

func unInstall() {
    root, err := util.GetRootDir()
    if err != nil {
		fmt.Println("获取根目录失败")
		os.Exit(1)
	}
    cmd.RemoveFromPath(root)
    fmt.Println("卸载成功")
}

func setJavaHome(jhome string) {
	config.Jhome = jhome
	cmd.SetEnvironmentValue("JAVA_HOME", jhome)
	writeConfig()
	fmt.Println("设置JAVA_HOME成功,需要重启终端生效")
}

func init() {
	// 加载配置文件
	loadConfig()
	// 初始化JAVA_HOME到Path
	cmd.AddToPath("%JAVA_HOME%\\bin")
}

func loadConfig() {
	root, err := util.GetRootDir()
	if err != nil {
		fmt.Println("获取根目录失败")
		os.Exit(1)
	}
	var configFile = path.Join(root, "jvm.yml")
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
		writeConfig()
	}
}


func help() {
	fmt.Println("jvm add <name> <path>               Add a JDK")
	fmt.Println("jvm rm <name>                       Remove a JDK")
	fmt.Println("jvm list                            List all installed JDKs")
	fmt.Println("jvm use <name>                      Use a JDK")
	fmt.Println("jvm home <path>                     Set the path of JAVA_HOME")
    fmt.Println("jvm global <install/uninstall>      Install jvm to system")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func writeConfig() {
	data, err := yaml.Marshal(config)
	if err != nil {
		fmt.Println("保存配置文件失败")
	}
	os.WriteFile(path.Join(config.Root, "jvm.yml"), data, 0644)
}

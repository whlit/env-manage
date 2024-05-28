package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
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

var client = &http.Client{}

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
	table.Printf()

}

// 添加 JDK
func add(version string, jpath string) {
	if !fileExists(jpath) {
		fmt.Println("路径不存在")
		return
	}
	if !fileExists(path.Join(jpath, "bin/java.exe")) {
		fmt.Println("路径不是 JDK 路径", jpath)
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
	var info DownloadInfo
	downloads := getDownloads()
	var options []huh.Option[DownloadInfo] = make([]huh.Option[DownloadInfo], len(downloads))
	for i, v := range downloads {
		options[i] = huh.NewOption(strings.Join([]string{v.Source, v.Version}, "_"), v)
	}
	huh.NewSelect[DownloadInfo]().Options(options...).Value(&info).Run()
	fmt.Println("开始安装JDK", info.Version)

	if !download(info) {
		fmt.Println("下载失败")
		return
	}

	// 下载完成 开始解压
	fmt.Println("正在解压...")
	zipfile := getDownloadFilePath(info)
	fileName := filepath.Base(zipfile)
	dir := filepath.Join(util.GetRootDir(), "versions", "jdk", strings.TrimSuffix(fileName, filepath.Ext(fileName)))
	if fileExists(dir) {
		os.RemoveAll(dir)
	}
	err := util.Unzip(zipfile, dir)
	if err != nil {
		fmt.Println("解压失败", err)
		return
	}
	// 解压完成 开始配置
	fmt.Println("解压完成, 正在添加到配置")
	dirs, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("读取目录失败", err)
		return
	}
	if len(dirs) == 1 {
		dir = filepath.Join(dir, dirs[0].Name())
	}
	add(strings.Join([]string{info.Source, info.Version}, "_"), dir)

	fmt.Println("安装成功")
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
	Source       string
}

func getDownloads() []DownloadInfo {
	return []DownloadInfo{
		{
			Version:      "jdk22",
			Url:          "https://download.oracle.com/java/22/latest/jdk-22_windows-x64_bin.zip",
			CheckType:    "sha256",
			CheckCodeUrl: "https://download.oracle.com/java/22/latest/jdk-22_windows-x64_bin.zip.sha256",
			Source:       "oracle",
		},
		{
			Version:      "jdk21",
			Url:          "https://download.oracle.com/java/21/latest/jdk-21_windows-x64_bin.zip",
			CheckType:    "sha256",
			CheckCodeUrl: "https://download.oracle.com/java/21/latest/jdk-21_windows-x64_bin.zip.sha256",
			Source:       "oracle",
		},
		{
			Version:      "jdk17",
			Url:          "https://download.oracle.com/java/17/latest/jdk-17_windows-x64_bin.zip",
			CheckType:    "sha256",
			CheckCodeUrl: "https://download.oracle.com/java/17/latest/jdk-17_windows-x64_bin.zip.sha256",
			Source:       "oracle",
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

func getDownloadFilePath(info DownloadInfo) string {
	fileName := filepath.Base(info.Url)
	// todo 从root 目录下创建一个临时目录
	// filePath := filepath.Join(util.GetRootDir(), "temp", fileName)
	filePath := filepath.Join("D:\\", "temp", fileName)
	return filePath
}

// 下载文件
// 下载文件到临时目录
func download(info DownloadInfo) bool {
	filePath := getDownloadFilePath(info)
	// 文件已经存在，说明下载过了，需要校验一下和最新的是否一致，不一致则删除，重新下载
	if fileExists(filePath) {
		ok, err := checkFile(info)
		if err != nil {
			fmt.Println("文件校验失败", err)
			os.Exit(1)
		}
		if ok {
			fmt.Println("版本未更新，不需要重新下载。")
			return true
		}
		os.Remove(filePath)
	} else {
		util.MkBaseDir(filePath)
	}
	// 创建文件
	out, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer out.Close()

	// 下载文件
	req, err := http.NewRequest("GET", info.Url, nil)
	if err != nil {
		fmt.Println(err)
		return false
	}
	req.Header.Set("User-Agent", "Env Manage")
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while downloading", "-", err)
		return false
	}
	defer response.Body.Close()
	fmt.Println("Downloading... ", info.Url)
	_, err = io.Copy(out, response.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("Download Completed")
	fmt.Println("校验文件...")
	// 下载完成 校验文件
	ok, err := checkFile(info)
	if err != nil {
		fmt.Println("文件校验失败", err)
		os.Exit(1)
	}
	return ok
}

// 校验文件
// 从远程获取校验码对下载的文件进行校验
func checkFile(info DownloadInfo) (bool, error) {
	filePath := getDownloadFilePath(info)
	if !fileExists(filePath) {
		return false, errors.New("文件不存在")
	}
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// 获取校验码
	req, err := http.NewRequest("GET", info.CheckCodeUrl, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("User-Agent", "Env Manage")
	response, err := client.Do(req)
	if err != nil {
		return false, err
	}
	code, err := io.ReadAll(response.Body)
	if err != nil {
		return false, err
	}

	// 计算校验码并进行校验
	switch info.CheckType {
	case "sha256":
		hash := sha256.New()
		_, err = io.Copy(hash, file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		hashInBytes := hash.Sum(nil)
		return hex.EncodeToString(hashInBytes) == string(code), nil
	}
	return true, nil
}

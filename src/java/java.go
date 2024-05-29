package java


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
	"github.com/whlit/env-manage/version"
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


// 列出所有已安装 JDK
func List() {
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
func Add(version string, jpath string) {
	if !util.FileExists(jpath) {
		fmt.Println("路径不存在")
		return
	}
	if !util.FileExists(path.Join(jpath, "bin/java.exe")) {
		fmt.Println("路径不是 JDK 路径", jpath)
	}
	config.Jdks[version] = jpath
	util.SaveConfig(config)
}

// 移除 JDK
func Remove(name string) {
	delete(config.Jdks, name)
	util.SaveConfig(config)
}

// 切换 JDK
func Use() {
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
func Home(jhomePath string) {
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

func Install() {
	var info version.VersionDownload
	downloads := version.GetJdkVersions()
	var options []huh.Option[version.VersionDownload] = make([]huh.Option[version.VersionDownload], len(downloads))
	for i, v := range downloads {
		options[i] = huh.NewOption(v.GetVersionKey(), v)
	}
	huh.NewSelect[version.VersionDownload]().Options(options...).Value(&info).Run()
	var confirm bool
	huh.NewConfirm().Title(strings.Join([]string{"确认安装 ", info.GetVersionKey(), " ?"}, "")).Value(&confirm).Run()
	if !confirm {
		fmt.Println("取消安装")
		return
	}
	fmt.Println("开始安装JDK", info.GetVersionKey())

	zipPath, err := info.Download()
	if err != nil {
		fmt.Println("下载失败")
		return
	}

	// 下载完成 开始解压
	fmt.Println("正在解压...")
	dir := filepath.Join(util.GetRootDir(), "versions", "jdk", info.GetVersionKey())
	if util.FileExists(dir) {
		os.RemoveAll(dir)
	}
	err = util.Unzip(zipPath, dir)
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
	Add(info.GetVersionKey(), dir)

	fmt.Println("安装成功")
}

func setJavaHome(jhome string) {
	config.Jhome = jhome
	cmd.SetEnvironmentValue("JAVA_HOME", jhome)
	util.SaveConfig(config)
	fmt.Println("设置JAVA_HOME成功,需要重启终端生效")
}

// 初始化
func InitConfig() {
	// 加载配置文件
	loadConfig()
	// 初始化JAVA_HOME到Path
	cmd.AddToPath("%JAVA_HOME%\\bin")
}

func loadConfig() {
	root := util.GetRootDir()
	var configFile = util.GetConfigFilePath()
	// 读取配置文件
	if util.FileExists(configFile) {
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
		config.Jhome = path.Join(root, "runtime", "jdk")
		setJavaHome(config.Jhome)
		util.MkBaseDir(config.Jhome)
		util.SaveConfig(config)
	}
}

func Help() {
	fmt.Println("add <name> <path>           Add a JDK")
	fmt.Println("rm <name>                   Remove a JDK")
	fmt.Println("list                        List all installed JDKs")
	fmt.Println("use                         Select And Use a JDK")
	fmt.Println("home <path>                 Set the path of JAVA_HOME")
}


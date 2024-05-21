package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"golang.org/x/text/encoding/simplifiedchinese"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Jdks map[string]string `yaml:"Jdks"`
	Jhome string           `yaml:"Jhome"`
	Root string            `yaml:"Root"`
}

var config = &Config{
	Jhome: filepath.Clean(os.Getenv("JAVA_HOME")),
	Jdks: make(map[string]string),
	Root: filepath.Clean(os.Getenv("ENV_ROOT")),
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
	case "root":
		root(args[2])
	default:
		help()
	}
}

func list() {
	for k, v := range config.Jdks {
		fmt.Println(k, v)
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
	home, _ := os.Lstat(config.Jhome)
	if home != nil {
		elevatedRun("rmdir", filepath.Clean(config.Jhome))
	}
	_, err := elevatedRun("mklink", "/D", filepath.Clean(config.Jhome), filepath.Clean(config.Jdks[name]))
	if err != nil {
		errr, _ := simplifiedchinese.GB18030.NewDecoder().String(err.Error())
		fmt.Println(errr)
	}
}

func root(rootPath string) {
	if fileExists(rootPath) {
		config.Root = rootPath
		writeConfig()
	}
}

func init() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("获取可执行文件失败")
		return
	}
	config.Root = filepath.Dir(exePath)
	var configFile = path.Join(filepath.Dir(exePath), "jvm.yml")
	if fileExists(configFile) {
		file, err := os.ReadFile(configFile)
		if err != nil {
			fmt.Println("读取配置文件失败")
		}
		var yamlData = &Config{}
		yaml.Unmarshal(file, &yamlData)
		if yamlData.Jdks != nil {
			config.Jdks = yamlData.Jdks
		}
		return
	}
	os.Create(configFile)
}

func elevatedRun(name string, arg ...string) (bool, error) {
	ok, err := run("cmd", nil, append([]string{"/C", name}, arg...)...)
	if err != nil {
		ok, err = run("elevate.cmd", &config.Root, append([]string{"cmd", "/C", name}, arg...)...)
	}
	return ok, err
}

func run(name string, dir *string, arg ...string) (bool, error) {
	c := exec.Command(name, arg...)
	if dir != nil {
		c.Dir = *dir
	}
	var stderr bytes.Buffer
	c.Stderr = &stderr
	err := c.Run()
	if err != nil {
		return false, errors.New(fmt.Sprint(err) + ": " + stderr.String())
	}

	return true, nil
}

func help() {
	fmt.Println("jvm add <name> <path>   Add a JDK")
	fmt.Println("jvm rm <name>           Remove a JDK")
	fmt.Println("jvm list                List all installed JDKs")
	fmt.Println("jvm use <name>          Use a JDK")
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

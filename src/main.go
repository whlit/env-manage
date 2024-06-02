package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/whlit/env-manage/core"
	"github.com/whlit/env-manage/java"
	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/maven"
	"github.com/whlit/env-manage/util"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		help()
		return
	}
	name := args[0]
	if len(args) < 2 {
		fmt.Println("需要指定要执行的动作: [add|rm|list|use|install|create]")
		help()
		return
	}
	action := args[1]
	switch name {
	case "jdk":
		config := core.NewConfig("jdk", &java.JdkEnvManager{
			EnvManager: core.EnvManager{
				EnvName:  "JAVA_HOME",
				Versions: make(map[string]string),
				Name:     "jdk",
			},
		})
		config.Load()
		manage(action, config.Data, args[2:])
		config.Save()
	case "maven":
		config := core.NewConfig("maven", &maven.MavenEnvManager{
			EnvManager: core.EnvManager{
				EnvName:  "M2_HOME",
				Versions: make(map[string]string),
				Name:     "maven",
			},
		})
		config.Load()
		manage(action, config.Data, args[2:])
		config.Save()
	default:
		config := core.NewConfig(name, &core.EnvManager{})
		config.Load()
		if config.Data == nil || config.Data.Name == "" {
			if action != "create" {
				logger.Error(fmt.Sprintf("不支持[%s]的版本管理,请先创建", name))
			}
			create(name)
			return
		}
		manage(action, config.Data, args[2:])
		config.Save()
	}
}

func manage(action string, manager core.IEnvManager, args []string) {
	switch action {
	case "add":
		if len(args) < 2 {
			logger.Error("Please provide a version")
		}
		manager.Add(args[0], args[1])
	case "rm":
		manager.Remove()
	case "use":
		manager.Use()
	case "list":
		manager.List()
	case "install":
		if _, ok := manager.(core.IInstallable); ok {
			manager.(core.IInstallable).Install()
		}
	}
}

func create(name string) {
	var envName string
	var envPathValue string
	huh.NewForm(huh.NewGroup(
		huh.NewInput().Title("输入要添加的环境变量名称:\n(例如:JDK环境变量名称为JAVA_HOME)").Value(&envName),
		huh.NewInput().Title("输入要添加到Path的值:\n(例如:JDK添加到Path为:%JAVA_HOME%\\bin").Value(&envPathValue),
	)).Run()
	envName = strings.TrimSpace(envName)
	envPathValue = strings.TrimSpace(envPathValue)
	if envName == "" || envPathValue == "" {
		logger.Error("输入不能为空")
	}

	envValue := filepath.Join(util.GetRootDir(), "runtime", name)
	table := util.NewTable("Name", "EnvName", "EnvValue", "EnvPathValue").Add(map[string]string{
		"Name":         name,
		"EnvName":      envName,
		"EnvValue":     envValue,
		"EnvPathValue": envPathValue,
	})
	var confirm bool
	huh.NewConfirm().Title(strings.Join(append([]string{"确认从以下信息创建版本管理?"}, table.Sprintf()...), "\n")).Value(&confirm).Run()
	if !confirm {
		logger.Error("用户取消操作")
	}
	config := core.NewConfig(name, &core.EnvManager{
		Name:         name,
		EnvName:      envName,
		EnvValue:     envValue,
		EnvPathValue: envPathValue,
		Versions:     make(map[string]string),
	})
	logger.Infof("创建版本管理:\n    Name: '%s',\n    EnvName: '%s',\n    EnvValue: '%s',\n    EnvPathValue: '%s'",
		config.Data.Name, config.Data.EnvName, config.Data.EnvValue, config.Data.EnvPathValue)
	config.Save()
	logger.Info("创建成功")
}

func help() {
	fmt.Println("Usage: vm <name> <action> [args]")
	fmt.Println("Name:                      环境管理的名称")
	fmt.Println("  jdk                      jdk版本管理")
	fmt.Println("  maven                    maven版本管理")
	fmt.Println("  <name>                   用create创建的其他版本管理的名称")
	fmt.Println("Actions:")
	fmt.Println("  create                   创建一个版本管理")
	fmt.Println("  add <version> <path>     添加版本,version: 版本名称(自定义),path: 版本的绝对路径")
	fmt.Println("  rm                       移除版本")
	fmt.Println("  list                     查询所有已添加的版本管理")
	fmt.Println("  use                      使用版本")
	fmt.Println("  install                  在线安装新版本,只支持jdk/maven的在线安装")
}

package main

import (
	"fmt"
	"os"

	"github.com/whlit/env-manage/core"
	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/manager"
)

var managers = make(map[string]core.IEnvManager)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		help()
		return
	}

	initManagers()

	if m, ok := managers[args[0]]; ok {
		manageEnv(args[1], m, args[2:])
	} else {
		help()
	}
}

func manageEnv(action string, manager core.IEnvManager, args []string) {
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
		manager.Install()
	default:
		help()
	}
}

func help() {
	fmt.Println("Usage: vm <name> <action> [args]")
	fmt.Println("Name:                      环境管理的名称")
	fmt.Println("  jdk                      jdk版本管理")
	fmt.Println("  maven                    maven版本管理")
	fmt.Println("  node                     node版本管理")
	fmt.Println("  <name>                   用create创建的其他版本管理的名称")
	fmt.Println("Actions:")
	fmt.Println("  create                   创建一个自定义版本管理")
	fmt.Println("  add <version> <path>     添加版本,version: 版本名称(自定义),path: 版本的绝对路径")
	fmt.Println("  rm                       移除版本")
	fmt.Println("  list                     查询所有已添加的版本管理")
	fmt.Println("  use                      使用版本")
	fmt.Println("  install                  在线安装新版本,自定义的版本管理不支持")
}

func initManagers() {
	for name, manager := range core.GlobalConfig.Managers {
		managers[name] = &manager
	}
	if m, ok := core.GlobalConfig.Managers["jdk"]; ok {
		managers["jdk"] = &manager.JdkEnvManager{EnvManager: m}
	} else {
		m := manager.NewJdkManager()
		managers["jdk"] = &manager.JdkEnvManager{EnvManager: m}
		core.GlobalConfig.Managers["jdk"] = m
		core.SaveConfig()
	}
}

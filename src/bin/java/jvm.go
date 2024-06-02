package main

import (
	"fmt"
	"os"

	"github.com/whlit/env-manage/core"
	"github.com/whlit/env-manage/java"
	"github.com/whlit/env-manage/logger"
)

var config = core.NewConfig("jdk", &java.JdkEnvManager{
	EnvManager: core.EnvManager{
		EnvName:  "JAVA_HOME",
		Versions: make(map[string]string),
		Name:     "jdk",
	},
})

func main() {
	args := os.Args
	action := ""
	if len(args) > 1 {
		action = args[1]
	}

	config.Load()

	switch action {
	case "list":
		config.Data.List()
	case "add":
		config.Data.Add(args[2], args[3])
	case "rm":
		config.Data.Remove()
	case "use":
		config.Data.Use()
	case "home":
		if config.Data.EnvValue == args[2] {
			return
		}
		_, err := config.Data.Home(args[2])
		if err != nil {
			logger.Error(err)
		}
		config.Save()
	case "install":
		config.Data.Install()
	default:
		help()
	}
}

func help() {
	fmt.Println("add <name> <path>           Add a JDK")
	fmt.Println("rm <name>                   Remove a JDK")
	fmt.Println("list                        List all installed JDKs")
	fmt.Println("use                         Select And Use a JDK")
	fmt.Println("home <path>                 Set the path of JAVA_HOME")
}

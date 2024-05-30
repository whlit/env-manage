package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/whlit/env-manage/java"
	"github.com/whlit/env-manage/maven"
	"github.com/whlit/env-manage/util"
)

func main() {
	args := os.Args
	action := ""
	if len(args) > 1 {
		action = args[1]
	}

	switch action {
	case "list":
		list(args[2:])
	case "install":
		install(args[2:])
	default:
		help()
	}
}

func list(args []string) {
	soft := getSoft(args)
	switch soft {
	case "jdk":
		java.List()
	case "maven":
		maven.List()
	default:
		fmt.Println("not support the value: ", soft)
	}
}

func install(args []string) {
	soft := getSoft(args)
	switch soft {
	case "jdk":
		java.Install()
	case "maven":
		maven.Install()
	default:
		fmt.Println("not support the value: ", soft)
	}
}

func help() {
	util.NewTable([]string{"Command", "Example", "Description"}).Add(
		map[string]string{
			"Command":     "list [jdk|maven]",
			"Description": "List all installed soft, If the parameter is empty, it will be selected manually",
			"Example":     "vm list jdk",
		},
		map[string]string{
			"Command":     "install [jdk|maven]",
			"Description": "Installed soft, If the parameter is empty, it will be selected manually",
			"Example":     "vm install jdk",
		},
	).Printf()
}

func getSoft(args []string) string {
	var soft string
	if len(args) > 0 {
		soft = args[0]
	} else {
		huh.NewSelect[string]().Height(5).Options(huh.NewOptions("jdk", "maven")...).Value(&soft).Run()
	}

	switch soft {
	case "jdk":
		java.InitConfig()
    case "maven":
        maven.InitConfig()
	}

	return soft
}

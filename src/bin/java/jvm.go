package main

import (
	"fmt"
	"os"

	"github.com/whlit/env-manage/java"
)

func main() {
	args := os.Args
	action := ""
	if len(args) > 1 {
		action = args[1]
	}

	java.InitConfig()

	switch action {
	case "list":
		java.List()
	case "add":
		java.Add(args[2], args[3])
	case "rm":
		java.Remove()
	case "use":
		java.Use()
	case "home":
		java.Home(args[2])
	case "install":
		java.Install()
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

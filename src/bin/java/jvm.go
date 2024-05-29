package main

import (
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
		java.Remove(args[2])
	case "use":
		java.Use()
	case "home":
		java.Home(args[2])
	case "install":
		java.Install()
	default:
		java.Help()
	}
}

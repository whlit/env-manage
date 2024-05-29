package main

import (
	"os"

	"github.com/whlit/env-manage/maven"
)

func main() {
	args := os.Args
	action := ""
	if len(args) > 1 {
		action = args[1]
	}

	switch action {
	case "list":
		maven.List()
	case "add":
		maven.Add(args[2], args[3])
	case "rm":
		maven.Remove(args[2])
	case "use":
		maven.Use()
	case "home":
		maven.Home(args[2])
	case "install":
		maven.Install()
	default:
		maven.Help()
	}
}

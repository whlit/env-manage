package main

import (
	"os"

	"github.com/whlit/env-manage/core"
	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/maven"
)

var config = core.NewConfig("maven", &maven.MavenEnvManager{
	EnvManager: core.EnvManager{
		EnvName:  "M2_HOME",
		Versions: make(map[string]string),
		Name:     "maven",
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

}

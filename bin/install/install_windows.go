package main

import (
	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/util"
)

func main() {
	root := util.GetExeDir()
    util.SetEnvs(map[string][]string{"VM_HOME": {root}})
    util.SetEnvs(map[string][]string{"PATH": {"%VM_HOME%\\bin"}})
	logger.Info("安装成功,请重新打开终端使用")
}

package main

import (
	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/util"
    windows_util "github.com/whlit/env-manage/bin/windows/util"
)

func main() {
	root := util.GetExeDir()
    windows_util.SetWindowsEnvs(map[string][]string{"VM_HOME": {root}})
    windows_util.SetWindowsEnvs(map[string][]string{"PATH": {"%VM_HOME%\\bin"}})
	logger.Info("安装成功,请重新打开终端使用")
}

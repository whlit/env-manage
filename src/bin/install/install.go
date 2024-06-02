package main

import (
	"github.com/whlit/env-manage/cmd"
	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/util"
)

func main() {
	root := util.GetExeDir()
	cmd.SetEnvironmentValue("VM_HOME", root)
	cmd.AddToPath("%VM_HOME%\\bin")
	logger.Info("安装成功,请重新打开终端使用")
}

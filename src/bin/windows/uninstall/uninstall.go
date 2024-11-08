package main

import (
	"github.com/whlit/env-manage/bin/windows/util"
	"github.com/whlit/env-manage/logger"
)

func main() {
    util.RemoveWindowsEnv("VM_HOME")
    util.RemoveFromPath("%VM_HOME%\\bin")
    logger.Info("卸载成功")
}

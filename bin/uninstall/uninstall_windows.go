package main

import (
	"github.com/whlit/env-manage/logger"
    "github.com/whlit/env-manage/util"
)

func main() {
    util.RemoveEnv("VM_HOME")
    util.RemoveFromPath("%VM_HOME%\\bin")
    logger.Info("卸载成功")
}

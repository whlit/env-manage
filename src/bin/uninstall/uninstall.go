package main

import (
	"github.com/whlit/env-manage/cmd"
	"github.com/whlit/env-manage/logger"
)

func main() {
	cmd.RemoveFromPath("%VM_HOME%\\bin")
    logger.Info("卸载成功,请重新打开终端使用")
}

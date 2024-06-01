package main

import (
	"fmt"

	"github.com/whlit/env-manage/cmd"
)

func main() {
	cmd.RemoveFromPath("%VM_HOME%\bin")
	fmt.Println("卸载成功,请重新打开终端使用")
}

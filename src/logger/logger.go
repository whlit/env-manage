package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/whlit/env-manage/util"
)

var logger *log.Logger


func init() {
	root, err := util.GetRootDir()
	if err != nil {
		fmt.Println("获取根目录失败", err)
		os.Exit(1)
	}
	file, err := os.OpenFile(filepath.Join(root, "sys.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalln("打开文件失败:", err)
    }
	logger = log.New(io.MultiWriter(file, os.Stderr), "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(format string, v ...any) {
	logger.Printf(format, v...)
}
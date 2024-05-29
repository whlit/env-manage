package logger

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var infoLogger *log.Logger
var errorLogger *log.Logger

func init() {
    exePath, err := os.Executable()
	if err != nil {
        fmt.Println("获取根目录失败:", err)
        os.Exit(1)
	}
	root := filepath.Dir(filepath.Dir(exePath))
	logFile := filepath.Join(root, "log", "sys.log")
	os.MkdirAll(filepath.Dir(logFile), fs.ModeDir)
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
        fmt.Println("打开日志文件失败:", err)
        os.Exit(1)
	}
	infoLogger = log.New(io.MultiWriter(file, os.Stderr), "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger = log.New(io.MultiWriter(file, os.Stderr), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(format string, v ...any) {
	infoLogger.Printf(format, v...)
}

func Error(v ...any) {
	errorLogger.Fatalln(v...)
}

func Warn(v ...any) {
	infoLogger.Panicln(v...)
}

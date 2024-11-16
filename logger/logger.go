package logger

import (
	"fmt"
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
	os.MkdirAll(filepath.Dir(logFile), 0755)
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("打开日志文件失败:", err)
		os.Exit(1)
	}
	infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Infof(format string, v ...any) {
	fmt.Printf(format, v...)
	fmt.Println()
	infoLogger.Printf(format, v...)
}

func Info(v ...any) {
	fmt.Println(v...)
	infoLogger.Println(v...)
}

func Error(v ...any) {
	fmt.Println(v...)
	errorLogger.Fatalln(v...)
}

func Warn(v ...any) {
	fmt.Println(v...)
	infoLogger.Panicln(v...)
}

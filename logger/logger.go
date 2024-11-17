package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	debugLogger *Logger
	infoLogger  *Logger
	warnLogger  *Logger
	errorLogger *Logger
	level       Level = INFO
)

const (
    DEBUG Level = iota
	INFO
	WARN
	ERROR
)

type Level int
type Logger struct {
	Level  Level
	Logger *log.Logger
}

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
	debugLogger = &Logger{Level: DEBUG, Logger: log.New(file, "DEBUG", log.Ldate|log.Ltime|log.Lshortfile)}
	infoLogger = &Logger{Level: INFO, Logger: log.New(file, "INFO", log.Ldate|log.Ltime|log.Lshortfile)}
	infoLogger = &Logger{Level: WARN, Logger: log.New(file, "WARN", log.Ldate|log.Ltime|log.Lshortfile)}
	infoLogger = &Logger{Level: ERROR, Logger: log.New(file, "ERROR", log.Ldate|log.Ltime|log.Lshortfile)}
}

func (l Level) String() string {
    switch l {
    case DEBUG:
        return "DEBUG"
    case INFO:
        return "INFO"
    case WARN:
        return "WARN"
    case ERROR:
        return "ERROR"
    default:
        return ""
    }
}

func SetLevel(l string) {
	switch l {
    case "DEBUG":
        level = DEBUG
    case "INFO":
        level = INFO
    case "WARN":
        level = WARN
    case "ERROR":
        level = ERROR
    }
}

func (l *Logger) Println(args ...any) {
	if level > l.Level {
		return
	}
	l.Logger.Println(args...)
}

func (l *Logger) Printf(format string, args ...any) {
	if level > l.Level {
		return
	}
	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}
	l.Logger.Printf(format, args...)
}



func Debug(v ...any) {
	debugLogger.Println(v...)
}

func Infof(format string, v ...any) {
	infoLogger.Printf(format, v...)
}

func Info(v ...any) {
	infoLogger.Println(v...)
}

func Error(v ...any) {
	errorLogger.Println(v...)
}

func Warn(v ...any) {
	warnLogger.Println(v...)
}

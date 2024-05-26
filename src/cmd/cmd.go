package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/util"
	"golang.org/x/sys/windows/registry"
)

// 设置用户环境变量
func SetUserEnvVar(name string, value string) {
	ElevatedRun("setx", name, value)
}

// 提升权限运行
func ElevatedRun(name string, arg ...string) (bool, error) {
	ok, err := Run("cmd", nil, append([]string{"/C", name}, arg...)...)
	if err != nil {
		root, e := util.GetRootDir()
		if e != nil {
			return false, err
		}
		ok, err = Run("elevate.cmd", &root, append([]string{"cmd", "/C", name}, arg...)...)
	}
	return ok, err
}

func CmdRun(name string, arg ...string) (bool, error) {
	return Run("cmd", nil, append([]string{"/C", name}, arg...)...)
}

func GetEnvironmentValue(name string) (string, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, "Environment", registry.QUERY_VALUE)
	if err != nil {
		fmt.Println("无法打开键：", err)
		return "", err
	}
	defer key.Close()
	value, _, err := key.GetStringValue(name)
	if err != nil {
		fmt.Println("无法读取值：", err)
		return "", err
	}
	return value, nil
}

func SetEnvironmentValue(name string, value string) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, "Environment", registry.QUERY_VALUE)
	if err != nil {
		fmt.Println("无法打开键：", err)
		return err
	}
	defer key.Close()
	oldValue, _, _ := key.GetStringValue(name)
	logger.Info("写入环境变量:%s, \n    旧值:'%s',  \n    新值:'%s'", name, oldValue, value)
	_, err = Run("reg", nil, "add", "HKEY_CURRENT_USER\\Environment", "/v", name, "/t", "REG_SZ", "/d", value, "/f")
	return err
}

func AddToPath(value string) {
	pathEnv, err := GetEnvironmentValue("Path")
	if err != nil {
		fmt.Println("获取Path环境变量失败")
	}
	paths := strings.Split(pathEnv, ";")
	var existed = false
	for _, item := range paths {
		if strings.Contains(item, value) {
			existed = true
		}
	}
	if !existed {
		var newPaths []string
		for _, item := range paths {
			if item != "" {
				newPaths = append(newPaths, item)
			}
		}
		if !existed {
			newPaths = append(newPaths, value)
		}
		SetEnvironmentValue("Path", strings.Join(newPaths, ";"))
	}
}

func RemoveFromPath(value string) {
	pathEnv, err := GetEnvironmentValue("Path")
	if err != nil {
		fmt.Println("获取Path环境变量失败")
	}
	paths := strings.Split(pathEnv, ";")
	var newPaths []string
	var existed = false
	for _, item := range paths {
		if item == value {
			existed = true
			continue
		}
		newPaths = append(newPaths, item)
	}
    if existed {
        SetEnvironmentValue("Path", strings.Join(newPaths, ";"))
    }
}

// 运行命令
func Run(name string, dir *string, arg ...string) (bool, error) {
	c := exec.Command(name, arg...)
	if dir != nil {
		c.Dir = *dir
	}
	var stderr bytes.Buffer
	c.Stderr = &stderr
	err := c.Run()
	if err != nil {
		return false, errors.New(fmt.Sprint(err) + ": " + stderr.String())
	}

	return true, nil
}

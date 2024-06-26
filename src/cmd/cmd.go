package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/util"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/text/encoding/simplifiedchinese"
)

// 获取注册表项
func GetEnvironmentValue(name string) (string, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, "Environment", registry.QUERY_VALUE)
	if err != nil {
		logger.Warn("无法打开键：", err)
		return "", err
	}
	defer key.Close()
	value, _, err := key.GetStringValue(name)
	if err != nil {
		logger.Warn("无法读取值：", err)
		return "", err
	}
	return value, nil
}

// 写入注册表项
func SetEnvironmentValue(name string, value string) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, "Environment", registry.QUERY_VALUE)
	if err != nil {
		logger.Warn("无法打开键：", err)
		return err
	}
	defer key.Close()
	oldValue, _, _ := key.GetStringValue(name)
	if oldValue == value {
		return nil
	}
	_, err = run("reg", nil, "add", "HKEY_CURRENT_USER\\Environment", "/v", name, "/t", "REG_SZ", "/d", value, "/f")
	logger.Infof("写入环境变量:%s, \n    旧值:'%s',  \n    新值:'%s'", name, oldValue, value)
	return err
}

// 向PATH添加目录
func AddToPath(value string) {
	pathEnv, err := GetEnvironmentValue("Path")
	if err != nil {
		logger.Error("获取Path环境变量失败", err)
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

// 从PATH中移除目录
func RemoveFromPath(value string) {
	pathEnv, err := GetEnvironmentValue("Path")
	if err != nil {
		logger.Error("获取Path环境变量失败", err)
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

func CreateLink(dir string, target string) {
	dir = filepath.Clean(dir)
	// 修改快捷方式指向
	home, _ := os.Lstat(dir)
	if home != nil {
		elevatedRun("rmdir", dir)
	}
	util.MkBaseDir(dir)
	_, err := elevatedRun("mklink", "/D", filepath.Clean(dir), target)
	if err != nil {
		errr, _ := simplifiedchinese.GB18030.NewDecoder().String(err.Error())
		logger.Error("创建快捷方式失败", errr)
	}
}

// 运行命令
func run(name string, dir *string, arg ...string) (bool, error) {
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

// 提升权限运行
func elevatedRun(name string, arg ...string) (bool, error) {
	ok, err := run("cmd", nil, append([]string{"/C", name}, arg...)...)
	if err != nil {
		root := util.GetRootDir()
		ok, err = run(".\\lib\\elevate.cmd", &root, append([]string{"cmd", "/C", name}, arg...)...)
	}
	return ok, err
}

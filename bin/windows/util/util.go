package util

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/whlit/env-manage/logger"
	"golang.org/x/sys/windows/registry"
)


func SetWindowsEnvs(envs map[string][]string) {
    for k, v := range envs {
        if k == "PATH" {
            pathValue, err := getRegValue("Path")
            if err != nil {
                logger.Warn("无法获取注册表值: Path ->", err)
            }
            items := appendEnvs(strings.Split(pathValue, ";"), v)
            setRegValue("Path", strings.Join(items, ";"))
            continue
        }
        setRegValue(k, strings.Join(v, ";"))
    }
}

func RemoveWindowsEnv(name string) error{
    oldValue, err := getRegValue(name)
    if err != nil {
        return err
    }
    if oldValue == "" {
        return nil
    }
    _, err = run("reg", nil, "delete", "HKEY_CURRENT_USER\\Environment", "/v", name, "/f")
    logger.Infof("删除环境变量:%s, \n    值:'%s'", name, oldValue)
    return err
}

// 获取Windows注册表项
func getRegValue(name string) (string, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, "Environment", registry.QUERY_VALUE)
	if err != nil {
		logger.Warn("无法打开键：", err)
		return "", err
	}
	defer key.Close()
	value, _, err := key.GetStringValue(name)
    if err == registry.ErrNotExist {
        return "", nil
    }
	if err != nil {
		logger.Warn("无法读取值：", err)
		return "", err
	}
	return value, nil
}

func setRegValue(name string, value string) error {
    oldValue, err := getRegValue(name)
    if err != nil {
        return err
    }
    if oldValue == value {
        return nil
    }
	_, err = run("reg", nil, "add", "HKEY_CURRENT_USER\\Environment", "/v", name, "/t", "REG_SZ", "/d", value, "/f")
	logger.Infof("写入环境变量:%s, \n    旧值:'%s',  \n    新值:'%s'", name, oldValue, value)
    return err
}


func appendEnvs(envs, items []string) []string {
    var res []string
    set := make(map[string]byte)
    for _, env := range envs {
        if env == "" {
            continue
        }
        set[env] = 0
        res = append(res, env)
    }
    for _, item := range items {
        if _, ok := set[item]; !ok {
            res = append(res, item)
        }
    }
    return res
}

// 创建文件夹链接
// linkPath 要创建的链接文件夹路径
// target 目标文件夹路径
// windows下使用目录链接 即 mklink /J linkPath target
func CreateLink(linkPath, target string) error {
    _, err := run("cmd", nil, "/C", "mklink", "/J", linkPath, target)
    return err
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


// 从PATH中移除目录
func RemoveFromPath(value string) {
	pathEnv, err := getRegValue("Path")
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
		setRegValue("Path", strings.Join(newPaths, ";"))
	}
}

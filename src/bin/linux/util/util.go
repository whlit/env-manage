package util

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/whlit/env-manage/logger"
	common_util "github.com/whlit/env-manage/util"
)

// 创建文件夹链接
// linkPath 要创建的链接文件夹路径
// target 目标文件夹路径
func CreateLink(linkPath, target string) error {
	return os.Symlink(target, linkPath)
}

func SetEnvs(envs map[string][]string) {
	oldEnvs := getOldEnvs()
	newEnvs := make([]string, 0)
	for _, oldEnv := range oldEnvs {
		if !strings.HasPrefix(oldEnv, "export") {
			newEnvs = append(newEnvs, oldEnv)
			continue
		}
		envItem := strings.Split(oldEnv, "=")
		oldEnvName := strings.Split(envItem[0], " ")[1]
		if _, ok := envs[oldEnvName]; !ok {
			newEnvs = append(newEnvs, oldEnv)
			continue
		}
		if oldEnvName == "PATH" {
			oldPaths := strings.Split(envItem[1], ":")
			newPaths := appendEnvs(oldPaths, envs["PATH"])
			newEnvs = append(newEnvs, "export PATH="+strings.Join(newPaths, ":"))
			delete(envs, "PATH")
			continue
		}
		newEnvs = append(newEnvs, fmt.Sprintf("export %s=%s", oldEnvName, strings.Join(envs[oldEnvName], ":")))
		delete(envs, oldEnvName)
	}

	for k, v := range envs {
		if k == "PATH" {
			newEnvs = append(newEnvs, fmt.Sprintf("export PATH=$PATH:%s", strings.Join(v, ":")))
			continue
		}
		newEnvs = append(newEnvs, fmt.Sprintf("export %s=%s", k, strings.Join(v, ":")))
	}
	err := os.WriteFile(filepath.Join(common_util.GetRootDir(), "config", "env.sh"), []byte(strings.Join(newEnvs, "\n")), 0644)
	if err != nil {
		logger.Error("写入文件失败：", err)
	}
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

func getOldEnvs() []string {
	envsPath := filepath.Join(common_util.GetRootDir(), "config", "env.sh")
	file, err := os.OpenFile(envsPath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.Error("打开文件失败：", err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	res := make([]string, 0)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			logger.Error("读取文件失败：", err)
		}
		res = append(res, string(line))
	}
	return res
}

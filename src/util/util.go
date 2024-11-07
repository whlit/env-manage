package util

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/whlit/env-manage/logger"
	"golang.org/x/sys/windows/registry"
)

var root string
var exeDir string
var client = &http.Client{}


// 获取根目录 获取失败则直接退出程序
// 本方法以当前可执行文件所在的目录为bin目录为前提
// 注意使用
func GetRootDir() string {
	if root != "" {
		return root
	}
	exePath, err := os.Executable()
	if err != nil {
		logger.Error("获取根目录失败", err)
	}
	// 软件目录为 bin 根目录应该为上级目录
	root = filepath.Dir(filepath.Dir(exePath))
	return root
}

// 获取当前可执行文件所在的目录
func GetExeDir() string {
	if exeDir != "" {
		return exeDir
	}
	exePath, err := os.Executable()
	if err != nil {
		logger.Error("获取可执行文件目录失败", err)
	}
	exeDir = filepath.Dir(exePath)
	return exeDir
}

// 获取下载文件夹
func GetDownloadDir() string {
	return filepath.Join(GetRootDir(), "download")
}

// 创建最后一个分隔符之前的目录
func MkBaseDir(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		_, err = os.Stat(filepath.Dir(path))
		if os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(path), fs.ModeDir)
		}
	}
}

// 解压缩
func Unzip(zipPath, dir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dir, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dir, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dir)+string(os.PathSeparator)) {
			logger.Error("illegal file path: ", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

// 文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// 创建文件夹链接
// linkPath 要创建的链接文件夹路径
// target 目标文件夹路径
// windows下使用目录链接 即 mklink /J linkPath target
func CreateLink(linkPath, target string) error {
    if runtime.GOOS == "windows" {
        _, err := run("cmd", nil, "/C", "mklink", "/J", linkPath, target)
        return err
    }
    return os.Symlink(linkPath, target)
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

func GetHttpClient() *http.Client {
    return client
}

// 发送get请求
func Get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Env Manage")
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	code, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return code, nil
}

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

// 获取Windows注册表项
func getRegValue(name string) (string, error) {
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

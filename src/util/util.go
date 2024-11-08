package util

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/whlit/env-manage/logger"
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
			os.MkdirAll(filepath.Dir(path), 00755)
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

var reg, _ = regexp.Compile(`[^\d|\.]`)

func CompareVersion(v1, v2 string) int {
	if v1 == v2 {
		return 0
	}
	v1 = reg.ReplaceAllString(v1, "")
	v2 = reg.ReplaceAllString(v2, "")
	v1s := strings.Split(v1, ".")
	v2s := strings.Split(v2, ".")
	l1 := len(v1s)
	l2 := len(v2s)
	for i := 0; i < l1 || i < l2; i++ {
		i1 := 0
		if i < l1 {
			i1, _ = strconv.Atoi(v1s[i])
		}
		i2 := 0
		if i < l2 {
			i2, _ = strconv.Atoi(v2s[i])
		}
		if i1 > i2 {
			return 1
		} else if i1 < i2 {
			return -1
		}
	}
	return 0
}

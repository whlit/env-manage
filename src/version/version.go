package version

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/util"
)

type VersionDownload interface {
	GetVersionKey() string
	Download() (string, error)
}

type EnvVersion struct {
	// 版本号
	Version string

	// 来源 厂商
	Source string

	// 下载链接
	Url string

	// 校验码
	CheckCode string

	// 校验类型
	CheckType string

	// 长期支持版本
	Lts bool

	// 是否最新版本
	Latest bool

	// 文件名称
	FileName string

	// 下载文件路径
	FilePath string
}

func (v *EnvVersion) GetVersionKey() string {
	return v.Source + "-" + v.Version
}

func (v *EnvVersion) Check(filePath string) (bool, error) {
	if !util.FileExists(filePath) {
		return false, errors.New("文件不存在")
	}
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// 计算校验码并进行校验
	switch v.CheckType {
	case "sha256":
		hash := sha256.New()
		_, err = io.Copy(hash, file)
		if err != nil {
			return false, err
		}
		hashInBytes := hash.Sum(nil)
		return hex.EncodeToString(hashInBytes) == v.CheckCode, nil
	}
	return true, nil
}

func (v *EnvVersion) Download(filePath string) bool {
	// 文件已经存在，说明下载过了，需要校验一下和最新的是否一致，不一致则删除，重新下载
	if util.FileExists(filePath) {
		ok, err := v.Check(filePath)
		if err != nil {
			logger.Error("校验文件失败：", err)
		}
		if ok {
			fmt.Println("版本未更新，不需要重新下载。")
			return true
		}
		os.Remove(filePath)
	} else {
		util.MkBaseDir(filePath)
	}
	// 创建文件
	out, err := os.Create(filePath)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer out.Close()

	// 下载文件
	req, err := http.NewRequest("GET", v.Url, nil)
	if err != nil {
		fmt.Println(err)
		return false
	}
	req.Header.Set("User-Agent", "Env Manage")
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while downloading", "-", err)
		return false
	}
	defer response.Body.Close()
	fmt.Println("Downloading... ", v.Url)
	_, err = io.Copy(out, response.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("Download Completed")
	fmt.Println("校验文件...")
	// 下载完成 校验文件
	ok, err := v.Check(filePath)
	if err != nil {
        logger.Error("校验文件失败：", err)
	}
	return ok
}

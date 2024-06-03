package version

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
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
}

// 获取版本标识
func (v *EnvVersion) GetVersionKey() string {
	return v.Source + "-" + v.Version
}

// 校验下载的文件
func (v *EnvVersion) Check(filePath string) (bool, error) {
    logger.Info("校验文件")
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
	case "sha512":
		hash := sha512.New()
		_, err = io.Copy(hash, file)
		if err != nil {
			return false, err
		}
		hashInBytes := hash.Sum(nil)
		return hex.EncodeToString(hashInBytes) == v.CheckCode, nil
	}
	return true, nil
}

// 下载文件到指定位置
func (v *EnvVersion) Download(filePath string) bool {
	// 文件已经存在，说明下载过了，需要校验一下和最新的是否一致，不一致则删除，重新下载
	if util.FileExists(filePath) {
		ok, err := v.Check(filePath)
		if err != nil {
			logger.Error("校验文件失败：", err)
		}
		if ok {
			logger.Info("版本未更新，不需要重新下载。")
			return true
		}
		os.Remove(filePath)
	} else {
		util.MkBaseDir(filePath)
	}
	// 创建文件
	out, err := os.Create(filePath)
	if err != nil {
		logger.Warn("创建文件失败：", err)
		return false
	}
	defer out.Close()

	// 下载文件
	req, err := http.NewRequest("GET", v.Url, nil)
	if err != nil {
		logger.Warn("创建请求失败：", err)
		return false
	}
	req.Header.Set("User-Agent", "Env Manage")
	response, err := client.Do(req)
	if err != nil {
		logger.Warn("下载失败：", err)
		return false
	}
	defer response.Body.Close()
	logger.Info("下载中...", v.Url)
	_, err = io.Copy(out, response.Body)
	if err != nil {
		logger.Warn("下载失败：", err)
		return false
	}
	logger.Info("下载完成")
	// 下载完成 校验文件
	ok, err := v.Check(filePath)
	if err != nil {
		logger.Error("校验文件失败：", err)
	}
	return ok
}

func Install(versions []VersionDownload, dir string) (VersionDownload, error) {
	if versions == nil || len(versions) < 1 {
		logger.Warn("未找到任何版本信息")
		return nil, errors.New("未找到任何版本信息")
	}
	var v VersionDownload
	var options []huh.Option[VersionDownload] = make([]huh.Option[VersionDownload], len(versions))
	for i, v := range versions {
		options[i] = huh.NewOption(v.GetVersionKey(), v)
	}
	huh.NewSelect[VersionDownload]().Options(options...).Value(&v).Run()
	var confirm bool
	huh.NewConfirm().Title(strings.Join([]string{"确认安装 ", v.GetVersionKey(), " ?"}, "")).Value(&confirm).Run()
	if !confirm {
		logger.Info("取消安装")
		os.Exit(1)
	}
	logger.Info("开始安装 ", v.GetVersionKey())

	zipPath, err := v.Download()
	if err != nil {
		logger.Warn("下载失败", err)
		return nil, err
	}

	// 下载完成 开始解压
	logger.Info("正在解压...")
	installPath := filepath.Join(dir, v.GetVersionKey())
	if util.FileExists(installPath) {
		os.RemoveAll(installPath)
	}
	err = util.Unzip(zipPath, installPath)
	if err != nil {
		logger.Warn("解压失败", err)
		return nil, err
	}
	// 解压完成 开始配置
	logger.Info("解压完成")

	return v, nil
}

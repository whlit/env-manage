package core

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/util"
)

type Version struct {
	Version  string `yaml:"version" json:"version"`
	Path     string `yaml:"path" json:"path"`
	Lts      bool   `yaml:"lts" json:"lts"`
	Os       string `yaml:"os" json:"os"`
	Arch     string `yaml:"arch" json:"arch"`
	Url      string `yaml:"url" json:"url"`
	Sum      string `yaml:"sum" json:"sum"`
	SumType  string `yaml:"sum_type" json:"sum_type"`
	Latest   bool   `yaml:"latest" json:"latest"`
	FileName string `yaml:"file_name" json:"file_name"`
	FileType string `yaml:"file_type" json:"file_type"`
	Size     int    `yaml:"size" json:"size"`
	App      string `yaml:"app" json:"app"`
}

func (v *Version) Download() error {
	path := v.GetDownloadFilePath()
	if util.FileExists(path) {
		if v.Check() {
			logger.Info("文件已存在，无需下载")
			return nil
		}
		os.Remove(path)
	}

    os.MkdirAll(filepath.Dir(path), fs.ModeDir)
	out, err := os.Create(path)
	if err != nil {
		logger.Error("创建文件失败：", err)
		return err
	}
	defer out.Close()

	if v.Url == "" {
		return errors.New("url is empty")
	}
	req, err := http.NewRequest("GET", v.Url, nil)
	if err != nil {
		logger.Error("创建请求失败：", err)
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36 Edg/130.0.0.0")
	response, err := util.GetHttpClient().Do(req)
	if err != nil {
		logger.Error("发送请求失败：", err)
		return err
	}
	defer response.Body.Close()

	logger.Info("下载中...", v.Url)
	_, err = io.Copy(out, response.Body)
	if err != nil {
		logger.Warn("下载失败：", err)
		return err
	}
	logger.Info("下载完成")

	if !v.Check() {
		logger.Warn("校验失败")
		return errors.New("文件校验失败")
	}

	logger.Info("校验成功")
	return nil
}

func (v *Version) Check() bool {
	path := v.GetDownloadFilePath()
	if !util.FileExists(path) {
		return false
	}

	if v.Sum == "" || v.SumType == "" {
		logger.Error("校验码为空")
		return false
	}

	file, err := os.Open(path)
	if err != nil {
		logger.Error("打开文件失败：", err)
		return false
	}
	defer file.Close()

	// 计算校验码并进行校验
	switch v.SumType {
	case "sha256":
		hash := sha256.New()
		_, err = io.Copy(hash, file)
		if err != nil {
			logger.Error("计算校验码失败：", err)
			return false
		}
		hashInBytes := hash.Sum(nil)
		return hex.EncodeToString(hashInBytes) == v.Sum
	case "sha512":
		hash := sha512.New()
		_, err = io.Copy(hash, file)
		if err != nil {
			logger.Error("计算校验码失败：", err)
			return false
		}
		hashInBytes := hash.Sum(nil)
		return hex.EncodeToString(hashInBytes) == v.Sum
	}
	return false
}

func (v *Version) GetDownloadFilePath() string {
	return filepath.Join(util.GetRootDir(), GlobalConfig.DownloadDir, v.App, v.FileName)
}

func (v *Version) GetVersionsPath() string {
	return filepath.Join(util.GetRootDir(), GlobalConfig.VersionsDir, v.App)
}

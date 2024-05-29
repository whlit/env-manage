package version

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/util"
)

var client = &http.Client{}

func GetJdkVersions() []VersionDownload {
	var res []VersionDownload
	res = append(res, getOracleVersions()...)
	res = append(res, getAdoptiumVersions()...)
	return res
}

func getDownloadFilePath(v EnvVersion) string {
	return filepath.Join(filepath.Join(util.GetRootDir(), "download", "jdk"), v.GetVersionKey(), filepath.Base(v.Url))
}

// OracleJdk
type OracleJdk struct {
	EnvVersion
}

// 获取OracleJDK的校验码
func (jdk OracleJdk) getCheckCode() string {
	code, err := get(jdk.EnvVersion.Url + ".sha256")
	if err != nil {
		logger.Error("获取OracleJDK校验码失败")
		os.Exit(1)
	}
	return string(code)
}

// 下载OracleJDK
func (jdk *OracleJdk) Download() (string, error) {
	filePath := getDownloadFilePath(jdk.EnvVersion)
	// 获取校验码
	jdk.EnvVersion.CheckCode = jdk.getCheckCode()
	jdk.EnvVersion.FilePath = filePath
	jdk.EnvVersion.FileName = filepath.Base(jdk.EnvVersion.FilePath)
	if jdk.EnvVersion.Download(filePath) {
		return filePath, nil
	}
	return filePath, errors.New("download failed")
}

// getOracleVersions 获取Oracle JDK 版本
func getOracleVersions() []VersionDownload {
	return []VersionDownload{
		&OracleJdk{
			EnvVersion: EnvVersion{
				Version:   "jdk22_latest",
				Url:       "https://download.oracle.com/java/22/latest/jdk-22_windows-x64_bin.zip",
				CheckType: "sha256",
				Source:    "oracle",
				Lts:       false,
				Latest:    true,
			},
		},
		&OracleJdk{
			EnvVersion: EnvVersion{
				Version:   "jdk21_latest",
				Url:       "https://download.oracle.com/java/21/latest/jdk-21_windows-x64_bin.zip",
				CheckType: "sha256",
				Source:    "oracle",
				Lts:       true,
				Latest:    true,
			},
		},
		&OracleJdk{
			EnvVersion: EnvVersion{
				Version:   "jdk17_latest",
				Url:       "https://download.oracle.com/java/17/latest/jdk-17_windows-x64_bin.zip",
				CheckType: "sha256",
				Source:    "oracle",
				Lts:       true,
				Latest:    true,
			},
		},
	}
}

type AdoptiumJdk struct {
	EnvVersion
	v int
}

func (jdk *AdoptiumJdk) Download() (string, error) {
	value, err := get(strings.Join([]string{"https://api.adoptium.net/v3/assets/latest/", strconv.Itoa(jdk.v), "/hotspot?architecture=x64&image_type=jdk&os=windows&vendor=eclipse"}, ""))
	if err != nil {
		return "", err
	}
	var assets []AdoptiumAssets
	json.Unmarshal(value, &assets)
	if assets == nil || len(assets) < 1 {
		return "", errors.New("未查询到JDK信息")
	}
	var asset AdoptiumAssets
	if len(assets) > 1 {
		var options []huh.Option[AdoptiumAssets] = make([]huh.Option[AdoptiumAssets], len(assets))
		for i, a := range assets {
			options[i] = huh.NewOption(a.ReleaseName, a)
		}
		huh.NewSelect[AdoptiumAssets]().Options(options...).Title("找到多个版本，请选择").Value(&asset).Run()
	} else {
		asset = assets[0]
	}

	jdk.EnvVersion.Url = asset.Binary.Package.Url
	filePath := getDownloadFilePath(jdk.EnvVersion)
	jdk.EnvVersion.FilePath = filePath
	jdk.EnvVersion.FileName = filepath.Base(jdk.EnvVersion.FilePath)
	jdk.EnvVersion.CheckCode = asset.Binary.Package.Checksum

	if jdk.EnvVersion.Download(filePath) {
		return filePath, nil
	}
	return filePath, errors.New("download failed")
}

type AdoptiumAssets struct {
	Binary struct {
		Package struct {
			Checksum string `json:"checksum"`
			Url      string `json:"link"`
		} `json:"package"`
	} `json:"binary"`
	ReleaseName string `json:"release_name"`
}

type AdoptiumAvailable struct {
	ReleasesLts              []int `json:"available_lts_releases"`
	Releases                 []int `json:"available_releases"`
	MostRecentFeatureRelease int   `json:"most_recent_feature_release"`
	MostRecentFeatureVersion int   `json:"most_recent_feature_version"`
	MostRecentLts            int   `json:"most_recent_lts"`
	TipVersion               int   `json:"tip_version"`
}

// 获取Adoptium版本
func getAdoptiumVersions() []VersionDownload {
	var versions []VersionDownload
	// 获取版本列表
	value, err := get("https://api.adoptium.net/v3/info/available_releases")
	if err != nil {
		return versions
	}
	adoptiumAvailable := AdoptiumAvailable{}
	json.Unmarshal(value, &adoptiumAvailable)

	// 构建版本列表
	for _, v := range adoptiumAvailable.Releases {
		var aj = AdoptiumJdk{
			EnvVersion: EnvVersion{
				Version:   strings.Join([]string{"jdk", strconv.Itoa(v), "_hotspot"}, ""),
				CheckType: "sha256",
				Source:    "adoptium",
				Latest:    true,
			},
			v: v,
		}
		for _, lts := range adoptiumAvailable.ReleasesLts {
			if v == lts {
				aj.EnvVersion.Lts = true
			}
		}
		versions = append(versions, &aj)
	}
	return versions
}

func get(url string) ([]byte, error) {
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

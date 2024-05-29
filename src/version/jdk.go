package version

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
)

var client = &http.Client{}

func GetJdkVersions() []VersionDownload {
	var res []VersionDownload
	oracleVersions := getOracleVersions()
	for _, v := range oracleVersions {
        v2 := v
		res = append(res, &v2)
	}
	// res = append(res, getAdoptiumVersions()...)
	return res
}

func getDownloadFilePath(v EnvVersion) string {
	// return filepath.Join(filepath.Join(util.GetRootDir(), "jdk"), v.GetVersionKey(), filepath.Base(v.Url))
	return filepath.Join("D:\\temp", v.GetVersionKey(), filepath.Base(v.Url))
}

type OracleJdk struct {
	EnvVersion
}

func (jdk OracleJdk) getCheckCode() string {
	code, err := get(jdk.EnvVersion.Url + ".sha256")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(code)
}

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
func getOracleVersions() []OracleJdk {
	return []OracleJdk{
		{
			EnvVersion: EnvVersion{
				Version:   "jdk22_latest",
				Url:       "https://download.oracle.com/java/22/latest/jdk-22_windows-x64_bin.zip",
				CheckType: "sha256",
				Source:    "oracle",
			},
		},
		{
			EnvVersion: EnvVersion{
				Version:   "jdk21_latest",
				Url:       "https://download.oracle.com/java/21/latest/jdk-21_windows-x64_bin.zip",
				CheckType: "sha256",
				Source:    "oracle",
			},
		},
		{
			EnvVersion: EnvVersion{
				Version:   "jdk17_latest",
				Url:       "https://download.oracle.com/java/17/latest/jdk-17_windows-x64_bin.zip",
				CheckType: "sha256",
				Source:    "oracle",
			},
		},
	}
}

// func getAdoptiumVersions() []EnvVersion {
// 	var versions []EnvVersion
// 	//https://mirrors.tuna.tsinghua.edu.cn/Adoptium/8/filelist
// 	// 获取版本列表
// 	value, err := get("https://mirrors.tuna.tsinghua.edu.cn/Adoptium/8/filelist")
// 	if err != nil {
// 		return versions
// 	}
// 	fmt.Println(string(value))
// 	return versions
// }

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

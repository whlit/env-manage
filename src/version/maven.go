package version

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/util"
)

func GetMavenVersions() []VersionDownload {
	return getApacheMavens()
}

type ApacheMaven struct {
	EnvVersion
}

func (maven *ApacheMaven) getDownloadFilePath() string {
    return filepath.Join(util.GetDownloadDir(), "maven", maven.GetVersionKey(), filepath.Base(maven.Url))
}

func (maven *ApacheMaven) getCheckCode() string {
	code, err := get(maven.EnvVersion.Url + ".sha512")
	if err != nil {
		logger.Error("获取ApacheMaven校验码失败")
		os.Exit(1)
	}
	return string(code)
}

func (maven *ApacheMaven) Download() (string, error) {
	filePath := maven.getDownloadFilePath()
	// 获取校验码
	maven.EnvVersion.CheckCode = maven.getCheckCode()
	if maven.EnvVersion.Download(filePath) {
		return filePath, nil
	}
	return filePath, errors.New("download failed")
}

func getApacheMavens() []VersionDownload {
	versions := []VersionDownload{}
	value, err := get("https://dlcdn.apache.org/maven/maven-3/")
	if err != nil {
		fmt.Println("获取Maven版本信息失败")
		return versions
	}
	re := regexp.MustCompile(`(>[\d\.]+/<)`)
	m := re.FindAllString(string(value), -1)

	urlFmt := "https://downloads.apache.org/maven/maven-3/%s/binaries/apache-maven-%s-bin.zip"
	for _, vstr := range m {
		vstr = vstr[1 : len(vstr)-2]
		versions = append(versions, &ApacheMaven{
			EnvVersion: EnvVersion{
				Version: vstr,
				Url:     fmt.Sprintf(urlFmt, vstr, vstr),
                Source:  "apache",
                CheckType: "sha512",
                Lts:      false,
                Latest:   false,
			},
		})
		fmt.Println(vstr)
	}
	return versions
}

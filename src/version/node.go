package version

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/whlit/env-manage/logger"
	"github.com/whlit/env-manage/util"
)

func GetNodeVersions() []VersionDownload {
	return getNodes()
}

type Node struct {
	EnvVersion
}

func (node *Node) GetVersionKey() string {
	return node.Version
}

func (node *Node) getDownloadFilePath() string {
	return filepath.Join(util.GetDownloadDir(), "node", node.GetVersionKey(), filepath.Base(node.Url))
}

func (maven *Node) getCheckCode() string {
    url := fmt.Sprintf("https://nodejs.org/dist/%s/SHASUMS256.txt", maven.Version)
    logger.Infof("获取校验码 %s", url)
	data, err := get(url)
	if err != nil {
		logger.Error("获取校验码失败")
		os.Exit(1)
	}
    scanner := bufio.NewScanner(bytes.NewReader(data))
    fileName := fmt.Sprintf("node-%s-win-x64.zip", maven.Version)
    for scanner.Scan() {
        line := scanner.Text()
        codes := strings.Split(line, "  ")
        if len(codes) >= 2 {
            if codes[1] == fileName {
                return codes[0]
            }
        }
    }
	return ""
}

func (maven *Node) Download() (string, error) {
	filePath := maven.getDownloadFilePath()
	// 获取校验码
	maven.EnvVersion.CheckCode = maven.getCheckCode()
	if maven.EnvVersion.Download(filePath) {
		return filePath, nil
	}
	return filePath, errors.New("download failed")
}

type NodeMsg struct {
	Version  string   `json:"version"`
	Files    []string `json:"files"`
	Npm      string   `json:"npm"`
	V8       string   `json:"v8"`
	Uv       string   `json:"uv"`
	Zlib     string   `json:"zlib"`
	Openssl  string   `json:"openssl"`
	Modules  string   `json:"modules"`
	Lts      bool     `json:"lts"`
	Security bool     `json:"security"`
}

func getNodes() []VersionDownload {
	versions := []VersionDownload{}
	value, err := get("https://nodejs.org/dist/index.json")
	if err != nil {
		logger.Warn("获取Maven版本信息失败", err)
		return versions
	}
	var msgs []NodeMsg
	json.Unmarshal(value, &msgs)

	for _, msg := range msgs {
		versions = append(versions, &Node{
			EnvVersion: EnvVersion{
                Version:   msg.Version,
                Source:    "nodejs",
				Url:       fmt.Sprintf("https://nodejs.org/dist/%s/node-%s-%s-x64.zip", msg.Version, msg.Version, "win"),
				CheckType: "sha256",
				Lts:       msg.Lts,
				Latest:    msg.Lts,
			},
		})
	}
	return versions
}

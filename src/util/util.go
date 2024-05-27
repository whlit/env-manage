package util

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// 获取根目录 获取失败则直接退出程序
// 本方法以当前可执行文件所在的目录为bin目录为前提
// 注意使用
func GetRootDir() string {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("获取根目录失败")
		os.Exit(1)
	}
	// 软件目录为 bin 根目录应该为上级目录
	return filepath.Dir(filepath.Dir(exePath))
}

// 获取当前可执行文件所在的目录
func GetExeDir() string {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("获取可执行文件目录失败")
		os.Exit(1)
	}
	return filepath.Dir(exePath)
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

// 获取配置文件路径
func GetConfigFilePath(name string) string {
	path := filepath.Join(GetRootDir(), "config", name)
	MkBaseDir(path)
	return path
}

type Table struct {
	Columns  []string
	Selected func(map[string]string) bool
	data     []map[string]string
	lens     map[string]int
}

func (t *Table) Print() {
	// 打印列名
	if t.Columns == nil {
		return
	}
	if t.Selected != nil {
		fmt.Print("   ")
	}
	var formats map[string]string = make(map[string]string)
	for _, column := range t.Columns {
		formats[column] = strings.Join([]string{"%-", strconv.Itoa(max(t.lens[column], len(column))), "s   "}, "")
		fmt.Printf(formats[column], column)
	}
	fmt.Print("\n")
	// 打印数据
	if t.data == nil {
		return
	}
	for _, row := range t.data {
		// 是否打印选中标记
		if t.Selected != nil {
			if t.Selected(row) {
				fmt.Print(" * ")
			} else {
				fmt.Print("   ")
			}
		}
		for _, column := range t.Columns {
			fmt.Printf(formats[column], row[column])
		}
		fmt.Print("\n")
	}
}

func (t *Table) Add(rows ...map[string]string) {
	if t.lens == nil {
		t.lens = make(map[string]int)
	}
	for _, row := range rows {
		for _, column := range t.Columns {
			t.lens[column] = max(t.lens[column], len(row[column]))
		}
	}
	t.data = append(t.data, rows...)
}

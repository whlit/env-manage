package main

import (
	"fmt"
	"path/filepath"

	"github.com/whlit/env-manage/cmd"
	"github.com/whlit/env-manage/util"
)

func main(){
    root := util.GetExeDir()
    cmd.AddToPath(filepath.Join(root, "bin"))
    fmt.Println("安装成功,请重新打开终端使用")
}

package main

import (
	tocfg "go2cfg/src"
	"go2cfg/src/config"
)

func main() {
	excelDir := config.Project.GetString("excel_dir")
	WritersStr := config.Project.Get("writers")
	tocfg.Start(tocfg.StartOptions{DirExcel: excelDir, WritersOptsStr: WritersStr})
}

package main

import (
	tocfg "go2cfg/src"
	"go2cfg/src/config"
)

func main() {
	excelDir := config.Project.GetString("excel_dir")
	outSvrDir := config.Project.GetString("out_svr_dir")
	outCliDir := config.Project.GetString("out_cli_dir")
	tocfg.Start(tocfg.StartOptions{DirExcel: excelDir, DirOutSvr: outSvrDir, DirOutCli: outCliDir})
}

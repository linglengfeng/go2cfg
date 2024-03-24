package main

import (
	tocfg "go2cfg/src"
	"go2cfg/src/config"
)

func main() {
	xlsxDir := config.Project.GetString("xlsx_dir")
	outDir := config.Project.GetString("out_dir")
	tocfg.Start(xlsxDir, outDir)
}

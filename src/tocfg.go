package tocfg

import (
	"fmt"
	"go2cfg/src/tojson"
	"os"
	"path/filepath"
	"strings"

	"github.com/tealeg/xlsx"
)

const (
	EXPORT_SVR      = "EXPORT_SVR"
	EXPORT_SVR_RAW  = 0
	EXPORT_SVR_CELL = 0
	EXPORT_SVR_IS   = 1
	EXPORT_SVR_NAME = 2

	EXPORT_CLI      = "EXPORT_CLI"
	EXPORT_CLI_RAW  = 1
	EXPORT_CLI_CELL = 0
	EXPORT_CLI_IS   = 1
	EXPORT_CLI_NAME = 2

	PRIMARY_KEY      = "PRIMARY_KEY"
	PRIMARY_KEY_RAW  = 2
	PRIMARY_KEY_CELL = 0
	PRIMARY_KEY_NAME = 1

	UNION_KEYS      = "UNION_KEYS"
	UNION_KEYS_RAW  = 3
	UNION_KEYS_CELL = 0
	UNION_KEYS_NAME = 1

	OUT_SVR      = "OUT_SVR"
	OUT_SVR_RAW  = 5
	OUT_SVR_CELL = 0

	OUT_CLI      = "OUT_CLI"
	OUT_CLI_RAW  = 6
	OUT_CLI_CELL = 0

	TYPE      = "TYPE"
	TYPE_RAW  = 7
	TYPE_CELL = 0

	NAME      = "NAME"
	NAME_RAW  = 8
	NAME_CELL = 0

	KEY      = "KEY"
	KEY_RAW  = 9
	KEY_CELL = 0

	NOTE      = "NOTE"
	NOTE_RAW  = 10
	NOTE_CELL = 0

	VALUE      = "VALUE"
	VALUE_CELL = 0

	TRUE  = "TRUE"
	FALSE = "FALSE"
	INT   = "INT"
	STR   = "STR"
	LIST  = "LIST"

	SVRDIR = "/svr"
	CLIDIR = "/cli"
)

type prop struct {
	Type string
	Raw  int
	Cell int
	Is   int
	Name int
	Val  any
}

type keyInfo struct {
	Key  string
	Cell int
}

type keyVal struct {
	Key any
	Val any
}

var (
	globalExportSvr  = prop{Type: EXPORT_SVR, Raw: EXPORT_SVR_RAW, Cell: EXPORT_SVR_CELL, Is: EXPORT_SVR_IS, Name: EXPORT_SVR_NAME}
	globalExportCli  = prop{Type: EXPORT_CLI, Raw: EXPORT_CLI_RAW, Cell: EXPORT_CLI_CELL, Is: EXPORT_CLI_IS, Name: EXPORT_CLI_NAME}
	globalPrimaryKey = prop{Type: PRIMARY_KEY, Raw: PRIMARY_KEY_RAW, Cell: PRIMARY_KEY_CELL, Name: PRIMARY_KEY_NAME}
	globalUnionKeys  = prop{Type: UNION_KEYS, Raw: UNION_KEYS_RAW, Cell: UNION_KEYS_CELL, Name: UNION_KEYS_NAME}
	globalOutSvr     = prop{Type: OUT_SVR, Raw: OUT_SVR_RAW, Cell: OUT_SVR_CELL}
	globalOutCli     = prop{Type: OUT_CLI, Raw: OUT_CLI_RAW, Cell: OUT_CLI_CELL}
	globalTypeProp   = prop{Type: TYPE, Raw: TYPE_RAW, Cell: TYPE_CELL}
	globalName       = prop{Type: NAME, Raw: NAME_RAW, Cell: NAME_CELL}
	globalKey        = prop{Type: KEY, Raw: KEY_RAW, Cell: KEY_CELL}
	globalNote       = prop{Type: NOTE, Raw: NOTE_RAW, Cell: NOTE_CELL}

	globalCheckList = []prop{globalExportSvr, globalExportCli, globalPrimaryKey, globalUnionKeys, globalOutSvr, globalOutCli, globalTypeProp, globalName, globalKey, globalNote}

	globalIntDefault  = 0
	globalStrDefault  = ""
	globalListDefault = []any{}

	globalCheckSameFile = map[string]struct{}{}

	globalExcelDir    = ""
	globalOutDir      = ""
	globalSvrFilePath = ""
	globalCliFilePath = ""

	globalTojsonWriter = new(tojson.Writer)
)

func Start(excelDir string, outDir string) {
	globalExcelDir = excelDir
	globalOutDir = outDir
	globalSvrFilePath = globalOutDir + SVRDIR
	globalCliFilePath = globalOutDir + CLIDIR

	files, err := os.ReadDir(globalExcelDir)
	if err != nil {
		fmt.Println("read dir error:", err)
	}
	createDirs()

	for _, file := range files {
		if file.Name()[0] == '$' || file.Name()[0] == '~' {
			continue
		}
		startSingle(file.Name())
	}
}

func startSingle(fileName string) {
	xlFile, err := xlsx.OpenFile(globalExcelDir + "/" + fileName)
	excelFileBase := filepath.Base(fileName)
	excelFileName := strings.TrimSuffix(excelFileBase, filepath.Ext(excelFileBase))
	if err != nil {
		fmt.Println("open Excel:", excelFileName, "err:", err)
		return
	}

	for _, sheet := range xlFile.Sheets {
		rowLen := len(sheet.Rows)
		if rowLen <= 0 {
			continue
		}
		if !checkRow(rowLen) {
			fmt.Println("file:", excelFileName, "sheet name:", sheet.Name, "row error, please check format.")
			return
		}
		if !checkCell(sheet.Rows) {
			fmt.Println("file:", excelFileName, "sheet name:", sheet.Name, "cell error, please check format.")
			return
		}
		setPropVal(excelFileName, sheet.Name, sheet.Rows)

		svrMapInfo := map[string]any{}
		cliMapInfo := map[string]any{}
		for _, row := range sheet.Rows {
			if len(row.Cells) <= 0 || row.Cells[VALUE_CELL].String() != VALUE {
				continue
			}

			addSvrMap := map[string]any{}
			addCliMap := map[string]any{}
			for cellIndex, cell := range row.Cells {
				text := cell.String()
				key := sheet.Rows[globalKey.Raw].Cells[cellIndex].String()
				if sheet.Rows[globalOutSvr.Raw].Cells[cellIndex].String() == TRUE {
					addSvrMap[key] = text
				}
				if sheet.Rows[globalOutCli.Raw].Cells[cellIndex].String() == TRUE {
					addCliMap[key] = text
				}
				// fmt.Printf("(%d,%d):%s\t", rowIndex, cellIndex, text)
				// fmt.Printf("%s\t", text)
			}
			// fmt.Println()
			idKey := row.Cells[VALUE_CELL+1].String()
			svrMapInfo[idKey] = addSvrMap
			cliMapInfo[idKey] = addSvrMap
		}

		globalTojsonWriter.SvrNameStr = fmt.Sprintf("%v", globalExportSvr.Val)
		globalTojsonWriter.CliNameStr = fmt.Sprintf("%v", globalExportCli.Val)
		globalTojsonWriter.SvrPath = globalSvrFilePath
		globalTojsonWriter.CliPath = globalCliFilePath
		globalTojsonWriter.SvrData = svrMapInfo
		globalTojsonWriter.CliData = cliMapInfo

		err = write()
		if err != nil {
			fmt.Println("file:", excelFileName, "sheet name:", sheet.Name, "write file error:", err)
			return
		}

		fmt.Println("file:", excelFileName, "sheet name:", sheet.Name, "ok.")
	}
}

func write() error {
	svrFile := globalTojsonWriter.SvrName()
	if _, ok := globalCheckSameFile[svrFile]; ok {
		return fmt.Errorf("filename repeat:%v", svrFile)
	}
	globalCheckSameFile[svrFile] = struct{}{}
	err := os.WriteFile(svrFile, append([]byte(globalTojsonWriter.ToSvrData()), byte('\n')), 0644)
	if err != nil {
		return err
	}

	cliFile := globalTojsonWriter.CliName()
	if _, ok := globalCheckSameFile[cliFile]; ok {
		return fmt.Errorf("filename repeat:%v", cliFile)
	}
	globalCheckSameFile[svrFile] = struct{}{}
	return os.WriteFile(cliFile, append([]byte(globalTojsonWriter.ToCliData()), byte('\n')), 0644)
}

func checkRow(len int) bool {
	rowLen := len - 1
	for _, prop := range globalCheckList {
		if rowLen < prop.Raw {
			return false
		}
	}
	return true
}

func checkCell(rows []*xlsx.Row) bool {
	for _, prop := range globalCheckList {
		if len(rows[prop.Raw].Cells)-1 < prop.Cell {
			return false
		}
	}
	return true
}

func setPropVal(fileName string, sheetName string, rows []*xlsx.Row) {
	for _, prop := range globalCheckList {
		switch prop.Type {
		case EXPORT_SVR:
			cells := rows[prop.Raw].Cells
			name := ""
			if len(cells)-1 >= prop.Name {
				name = fmt.Sprintf("%v", rows[prop.Raw].Cells[prop.Name])
			} else {
				name = fileName + "_" + sheetName
			}
			globalExportSvr.Val = name
		case EXPORT_CLI:
			cells := rows[prop.Raw].Cells
			name := ""
			if len(cells)-1 >= prop.Name {
				name = fmt.Sprintf("%v", rows[prop.Raw].Cells[prop.Name])
			} else {
				name = fileName + "_" + sheetName
			}
			globalExportCli.Val = name
		default:
			continue
		}
	}
}

func createDirs() {
	createDir(globalSvrFilePath)
	createDir(globalCliFilePath)
}

func createDir(directoryPath string) {
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		err := os.MkdirAll(directoryPath, 0644)
		if err != nil {
			fmt.Println("create Dir:", directoryPath, "err:", err)
		}
	}
}

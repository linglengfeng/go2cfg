package tocfg

import (
	"encoding/json"
	"fmt"
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
	exportSvr  = prop{Type: EXPORT_SVR, Raw: EXPORT_SVR_RAW, Cell: EXPORT_SVR_CELL, Is: EXPORT_SVR_IS, Name: EXPORT_SVR_NAME}
	exportCli  = prop{Type: EXPORT_CLI, Raw: EXPORT_CLI_RAW, Cell: EXPORT_CLI_CELL, Is: EXPORT_CLI_IS, Name: EXPORT_CLI_NAME}
	primaryKey = prop{Type: PRIMARY_KEY, Raw: PRIMARY_KEY_RAW, Cell: PRIMARY_KEY_CELL, Name: PRIMARY_KEY_NAME}
	unionKeys  = prop{Type: UNION_KEYS, Raw: UNION_KEYS_RAW, Cell: UNION_KEYS_CELL, Name: UNION_KEYS_NAME}
	outSvr     = prop{Type: OUT_SVR, Raw: OUT_SVR_RAW, Cell: OUT_SVR_CELL}
	outCli     = prop{Type: OUT_CLI, Raw: OUT_CLI_RAW, Cell: OUT_CLI_CELL}
	typeProp   = prop{Type: TYPE, Raw: TYPE_RAW, Cell: TYPE_CELL}
	name       = prop{Type: NAME, Raw: NAME_RAW, Cell: NAME_CELL}
	key        = prop{Type: KEY, Raw: KEY_RAW, Cell: KEY_CELL}
	note       = prop{Type: NOTE, Raw: NOTE_RAW, Cell: NOTE_CELL}

	checkList = []prop{exportSvr, exportCli, primaryKey, unionKeys, outSvr, outCli, typeProp, name, key, note}

	intDefault  = 0
	strDefault  = ""
	listDefault = []any{}

	checkSameFile = map[string]struct{}{}
)

func Start(wordDir string, outDir string) {
	files, err := os.ReadDir(wordDir)
	if err != nil {
		fmt.Println("read dir error:", err)
	}

	for _, file := range files {
		if file.Name()[0] == '$' || file.Name()[0] == '~' {
			continue
		}
		startSingle(wordDir, outDir, file.Name())
	}
}

func startSingle(wordDir string, outDir string, fileName string) {
	createDir(outDir + SVRDIR)
	createDir(outDir + CLIDIR)
	xlFile, err := xlsx.OpenFile(wordDir + "/" + fileName)
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
		setPropVal(outDir, excelFileName, sheet.Name, sheet.Rows)
		// svrKeys, cliKeys := keyCell(sheet.Rows)
		// svrInfo := map[string]any{}
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
				key := sheet.Rows[key.Raw].Cells[cellIndex].String()
				if sheet.Rows[outSvr.Raw].Cells[cellIndex].String() == TRUE {
					addSvrMap[key] = text
				}
				if sheet.Rows[outCli.Raw].Cells[cellIndex].String() == TRUE {
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

		writeSvrFile := fmt.Sprintf("%v", exportSvr.Val)
		if sheet.Rows[exportSvr.Raw].Cells[exportSvr.Is].String() == TRUE {
			err = write(writeSvrFile, svrMapInfo)
			if err != nil {
				fmt.Println("file:", excelFileName, "sheet name:", sheet.Name, "out svr file:", writeSvrFile, "error:", err)
				return
			}
		}

		if sheet.Rows[exportCli.Raw].Cells[exportCli.Is].String() == TRUE {
			writeCliFile := fmt.Sprintf("%v", exportCli.Val)
			err = write(writeCliFile, svrMapInfo)
			if err != nil {
				fmt.Println("file:", excelFileName, "sheet name:", sheet.Name, "out cli file:", writeCliFile, "error:", err)
				return
			}
		}

		fmt.Println("file:", excelFileName, "sheet name:", sheet.Name, "ok.")
	}
}

func write(outFile string, mapInfo map[string]any) error {
	if _, ok := checkSameFile[outFile]; ok {
		return fmt.Errorf("filename repeat:%v", outFile)
	}
	checkSameFile[outFile] = struct{}{}
	content := map2Json(mapInfo)
	err := os.WriteFile(outFile, append([]byte(content), byte('\n')), 0644)
	return err
}

func checkRow(len int) bool {
	rowLen := len - 1
	for _, prop := range checkList {
		if rowLen < prop.Raw {
			return false
		}
	}
	return true
}

func checkCell(rows []*xlsx.Row) bool {
	for _, prop := range checkList {
		if len(rows[prop.Raw].Cells)-1 < prop.Cell {
			return false
		}
	}
	return true
}

func setPropVal(outDir string, fileName string, sheetName string, rows []*xlsx.Row) {
	for _, prop := range checkList {
		switch prop.Type {
		case EXPORT_SVR:
			cells := rows[prop.Raw].Cells
			name := ""
			if len(cells)-1 >= prop.Name {
				name = fmt.Sprintf("%v", rows[prop.Raw].Cells[prop.Name])
			} else {
				name = fileName + "_" + sheetName
			}
			exportSvr.Val = outDir + SVRDIR + "/" + name + ".json"
		case EXPORT_CLI:
			cells := rows[prop.Raw].Cells
			name := ""
			if len(cells)-1 >= prop.Name {
				name = fmt.Sprintf("%v", rows[prop.Raw].Cells[prop.Name])
			} else {
				name = fileName + "_" + sheetName
			}
			exportCli.Val = outDir + CLIDIR + "/" + name + ".json"
		default:
			continue
		}
	}
}

// func keyCell(Rows []*xlsx.Row) ([]keyInfo, []keyInfo) {
// 	svrKeys := []keyInfo{}
// 	cliKeys := []keyInfo{}
// 	for i, cell := range Rows[key.Raw].Cells {
// 		if Rows[outSvr.Raw].Cells[i].String() == TRUE {
// 			svrKeys = append(svrKeys, keyInfo{Key: cell.String(), Cell: i})
// 		}
// 		if Rows[outCli.Raw].Cells[i].String() == TRUE {
// 			svrKeys = append(cliKeys, keyInfo{Key: cell.String(), Cell: i})
// 		}
// 	}
// 	return svrKeys, cliKeys
// }

func map2Json(param map[string]interface{}) string {
	dataType, err := json.Marshal(param)
	if err != nil {
		panic(err)
	}
	dataString := string(dataType)
	return dataString
}

func createDir(directoryPath string) {
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		err := os.MkdirAll(directoryPath, 0644)
		if err != nil {
			fmt.Println("create Dir:", directoryPath, "err:", err)
		}
	}
}

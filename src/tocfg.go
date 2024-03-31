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

type PrimarykeyVal struct {
	key valInfo
	val []valInfo
}

type valInfo struct {
	Name string
	Type string
	Val  string
}

type UnionKeysVals struct {
	Keys []valInfo
	Vals [][]valInfo
}

var (
	globalExportSvr  = prop{Type: EXPORT_SVR, Raw: EXPORT_SVR_RAW, Cell: EXPORT_SVR_CELL, Is: EXPORT_SVR_IS, Name: EXPORT_SVR_NAME}
	globalExportCli  = prop{Type: EXPORT_CLI, Raw: EXPORT_CLI_RAW, Cell: EXPORT_CLI_CELL, Is: EXPORT_CLI_IS, Name: EXPORT_CLI_NAME}
	globalPrimaryKey = prop{Type: PRIMARY_KEY, Raw: PRIMARY_KEY_RAW, Cell: PRIMARY_KEY_CELL, Name: PRIMARY_KEY_NAME}
	globalUnionKeys  = prop{Type: UNION_KEYS, Raw: UNION_KEYS_RAW, Cell: UNION_KEYS_CELL, Name: UNION_KEYS_NAME}
	globalOutSvr     = prop{Type: OUT_SVR, Raw: OUT_SVR_RAW, Cell: OUT_SVR_CELL}
	globalOutCli     = prop{Type: OUT_CLI, Raw: OUT_CLI_RAW, Cell: OUT_CLI_CELL}
	globalType       = prop{Type: TYPE, Raw: TYPE_RAW, Cell: TYPE_CELL}
	globalName       = prop{Type: NAME, Raw: NAME_RAW, Cell: NAME_CELL}
	globalKey        = prop{Type: KEY, Raw: KEY_RAW, Cell: KEY_CELL}
	globalNote       = prop{Type: NOTE, Raw: NOTE_RAW, Cell: NOTE_CELL}

	globalCheckList = []prop{globalExportSvr, globalExportCli, globalPrimaryKey, globalUnionKeys, globalOutSvr, globalOutCli, globalType, globalName, globalKey, globalNote}

	globalIntDefault  = 0
	globalStrDefault  = ""
	globalListDefault = []any{}

	globalCheckSameFile = map[string]struct{}{}

	globalExcelDir    = ""
	globalOutDir      = ""
	globalSvrFilePath = ""
	globalCliFilePath = ""

	globalPrimaryKeyInfo = valInfo{}
	globalUnionKyesInfo  = []valInfo{}

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
		svrSliceInfo := []PrimarykeyVal{}
		cliSliceInfo := []PrimarykeyVal{}
		for rowIndex, row := range sheet.Rows {
			if rowIndex == globalKey.Raw {
				for cellIndex, cell := range row.Cells {
					text := cell.String()
					if text == globalPrimaryKey.Val {
						globalPrimaryKeyInfo.Name = text
						globalPrimaryKeyInfo.Val = text
						globalPrimaryKeyInfo.Type = sheet.Rows[globalType.Raw].Cells[cellIndex].String()
					}

				}
			}
			if len(row.Cells) <= 0 || row.Cells[VALUE_CELL].String() != VALUE {
				continue
			}

			addSvrMap := map[string]any{}
			addCliMap := map[string]any{}
			addSvrSlice := []valInfo{}
			addCliSlice := []valInfo{}
			for cellIndex, cell := range row.Cells {
				text := cell.String()
				addkeyName := sheet.Rows[globalKey.Raw].Cells[cellIndex].String()
				addkeyType := sheet.Rows[globalType.Raw].Cells[cellIndex].String()
				if sheet.Rows[globalOutSvr.Raw].Cells[cellIndex].String() == TRUE {
					addSvrMap[addkeyName] = text
					addSvrSlice = append(addSvrSlice, valInfo{Type: addkeyType, Val: text, Name: addkeyName})
				}
				if sheet.Rows[globalOutCli.Raw].Cells[cellIndex].String() == TRUE {
					addCliMap[addkeyName] = text
					addCliSlice = append(addCliSlice, valInfo{Type: addkeyType, Val: text, Name: addkeyName})
				}
				// fmt.Printf("(%d,%d):%s\t", rowIndex, cellIndex, text)
				// fmt.Printf("%s\t", text)
			}
			// fmt.Println()
			idVal := row.Cells[VALUE_CELL+1].String()
			idType := sheet.Rows[globalType.Raw].Cells[VALUE_CELL+1].String()
			idName := sheet.Rows[globalName.Raw].Cells[VALUE_CELL+1].String()
			// svrMapInfo[idKey] = addSvrMap
			// cliMapInfo[idKey] = addSvrMap
			svrSliceInfo = append(svrSliceInfo, PrimarykeyVal{key: valInfo{Type: idType, Val: idVal, Name: idName}, val: addSvrSlice})
			cliSliceInfo = append(cliSliceInfo, PrimarykeyVal{key: valInfo{Type: idType, Val: idVal, Name: idName}, val: addCliSlice})
		}
		// globalPrimaryKey
		// fmt.Println("====================== globalPrimaryKey:", globalPrimaryKey)
		// fmt.Println("====================== globalUnionKeys:", globalUnionKeys)
		// fmt.Println("====================== globalPrimaryKeyInfo:", globalPrimaryKeyInfo)
		// fmt.Println("====================== globalUnionKeyInfo:", globalUnionKyesInfo)
		// fmt.Println("====================== svrSliceInfo:", svrSliceInfo)
		// fmt.Println("====================== cliSliceInfo:", cliSliceInfo)

		globalTojsonWriter.SvrNameStr = fmt.Sprintf("%v", globalExportSvr.Val)
		globalTojsonWriter.CliNameStr = fmt.Sprintf("%v", globalExportCli.Val)
		globalTojsonWriter.SvrPath = globalSvrFilePath
		globalTojsonWriter.CliPath = globalCliFilePath
		globalTojsonWriter.SvrData = svrMapInfo
		globalTojsonWriter.CliData = cliMapInfo

		//  []map[string][][]valInfo
		svrMapListInfo, cliMapListInfo := generateUnionKeysInfo(svrSliceInfo, cliSliceInfo)
		err = writeMap(svrSliceInfo, cliSliceInfo, svrMapListInfo, cliMapListInfo)
		// err = write()
		if err != nil {
			fmt.Println("file:", excelFileName, "sheet name:", sheet.Name, "write file error:", err)
			return
		}

		fmt.Println("file:", excelFileName, "sheet name:", sheet.Name, "ok.")
	}
}

func writeMap(svrSliceInfo []PrimarykeyVal, cliSliceInfo []PrimarykeyVal, svrMapListInfo []map[string][][]valInfo, cliMapListInfo []map[string][][]valInfo) error {
	svrFile := globalTojsonWriter.SvrName()
	if _, ok := globalCheckSameFile[svrFile]; ok {
		return fmt.Errorf("filename repeat:%v", svrFile)
	}
	globalCheckSameFile[svrFile] = struct{}{}
	// tostring
	svrdataPrimarykeyJsonStr := PrimarykeyJsonStr(svrSliceInfo)
	// svrdataUnionKeysJsonStr := UnionKeysJsonStr(svrMapListInfo) // todo
	svrdata := svrdataPrimarykeyJsonStr
	err := os.WriteFile(svrFile, append([]byte(svrdata), byte('\n')), 0644)
	if err != nil {
		return err
	}

	// cli----------
	cliFile := globalTojsonWriter.CliName()
	if _, ok := globalCheckSameFile[cliFile]; ok {
		return fmt.Errorf("filename repeat:%v", cliFile)
	}
	globalCheckSameFile[cliFile] = struct{}{}
	//tostring
	clidataPrimarykeyJsonStr := PrimarykeyJsonStr(cliSliceInfo)
	clidata := clidataPrimarykeyJsonStr
	err = os.WriteFile(cliFile, append([]byte(clidata), byte('\n')), 0644)
	if err != nil {
		return err
	}
	return nil
}

func UnionKeysJsonStr(svrMapListInfo []map[string][][]valInfo) string {

	return ""
}

func PrimarykeyJsonStr(svrSliceInfo []PrimarykeyVal) string {
	data := "{\n"
	for kvi, kv := range svrSliceInfo {
		kvstring := "\t\"" + kv.key.Val + "\"" + ": {\n"
		vstring := ""

		for iv, v := range kv.val {
			vvalstr := ""
			switch v.Type {
			case STR:
				vvalstr = "\"" + v.Val + "\""
				// fmt.Println("type:", v.Type, v.Val, vvalstr)
			default:
				vvalstr = v.Val
			}
			if len(kv.val)-1 == iv {
				vstring += "\t\t\"" + v.Name + "\"" + ": " + vvalstr + "\n"
			} else {
				vstring += "\t\t\"" + v.Name + "\"" + ": " + vvalstr + ",\n"
			}
		}
		vstring += "\n"
		if len(svrSliceInfo)-1 == kvi {
			kvstring += vstring + "\t}\n"
		} else {
			kvstring += vstring + "\t},\n"
		}

		data += kvstring
	}
	data += "}"
	return data
}

// func write() error {
// 	svrFile := globalTojsonWriter.SvrName()
// 	if _, ok := globalCheckSameFile[svrFile]; ok {
// 		return fmt.Errorf("filename repeat:%v", svrFile)
// 	}
// 	globalCheckSameFile[svrFile] = struct{}{}
// 	err := os.WriteFile(svrFile, append([]byte(globalTojsonWriter.ToSvrData()), byte('\n')), 0644)
// 	if err != nil {
// 		return err
// 	}

// 	cliFile := globalTojsonWriter.CliName()
// 	if _, ok := globalCheckSameFile[cliFile]; ok {
// 		return fmt.Errorf("filename repeat:%v", cliFile)
// 	}
// 	globalCheckSameFile[svrFile] = struct{}{}
// 	return os.WriteFile(cliFile, append([]byte(globalTojsonWriter.ToCliData()), byte('\n')), 0644)
// }

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
				name = rows[prop.Raw].Cells[prop.Name].String()
			} else {
				name = fileName + "_" + sheetName
			}
			globalExportSvr.Val = name
		case EXPORT_CLI:
			cells := rows[prop.Raw].Cells
			name := ""
			if len(cells)-1 >= prop.Name {
				name = rows[prop.Raw].Cells[prop.Name].String()
			} else {
				name = fileName + "_" + sheetName
			}
			globalExportCli.Val = name
		case PRIMARY_KEY:
			globalPrimaryKey.Val = rows[prop.Raw].Cells[prop.Name].String()
		case UNION_KEYS:
			globalUnionKeys.Val = rows[prop.Raw].Cells[prop.Name].String()
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

func generateUnionKeysInfo(svrSliceInfo []PrimarykeyVal, cliSliceInfo []PrimarykeyVal) ([]map[string][][]valInfo, []map[string][][]valInfo) {
	unionKeys := splitKeys(fmt.Sprintf("%v", globalUnionKeys.Val))
	// fmt.Println("====================== unionKeys:", unionKeys)
	// unionSingleMap := make(map[string][][]valInfo)              // 单个map
	// unionInfo := make([]map[string]([][]valInfo), len(unionKeys)) //所有联合key的 map集合
	unionSvrInfo := make([]map[string]([][]valInfo), len(unionKeys))
	for i := range unionKeys {
		unionSvrInfo[i] = make(map[string]([][]valInfo))
	}
	for _, svrSliceInfoSingle := range svrSliceInfo {
		singleUkeys := makeSingleUkeys(unionKeys)       // 五个key ，，val 都是 svrSliceInfoSingle.val
		singleUkeysStr := makeSingleUkeysStr(unionKeys) // 方面用map，key直接转成string
		for _, singleVal := range svrSliceInfoSingle.val {
			for unionKeyIndex, unionKey := range unionKeys {
				isin, inUnionpos := inUnionKey(singleVal.Name, unionKey)
				if isin {
					singleUkeys[unionKeyIndex][inUnionpos] = singleVal
					singleUkeysStr[unionKeyIndex][inUnionpos] = valInfo2Str(singleVal)

				}
			}
		}
		// keys 生成结束
		for i, slist := range singleUkeysStr {
			unionSingleMapVal := unionSvrInfo[i]
			keysOneStr := keys2Str(slist)
			unionSingleMapVal[keysOneStr] = append(unionSingleMapVal[keysOneStr], svrSliceInfoSingle.val)
			unionSvrInfo[i] = unionSingleMapVal
		}
	}

	// client
	unionCliInfo := make([]map[string]([][]valInfo), len(unionKeys))
	for i := range unionKeys {
		unionCliInfo[i] = make(map[string]([][]valInfo))
	}
	for _, cliSliceInfoSingle := range cliSliceInfo {
		singleUkeys := makeSingleUkeys(unionKeys)       // 五个key ，，val 都是 svrSliceInfoSingle.val
		singleUkeysStr := makeSingleUkeysStr(unionKeys) // 方面用map，key直接转成string
		for _, singleVal := range cliSliceInfoSingle.val {
			for unionKeyIndex, unionKey := range unionKeys {
				isin, inUnionpos := inUnionKey(singleVal.Name, unionKey)
				if isin {
					singleUkeys[unionKeyIndex][inUnionpos] = singleVal
					singleUkeysStr[unionKeyIndex][inUnionpos] = valInfo2Str(singleVal)

				}
			}
		}
		// keys 生成结束
		for i, slist := range singleUkeysStr {
			unionSingleMapVal := unionCliInfo[i]
			keysOneStr := keys2Str(slist)
			unionSingleMapVal[keysOneStr] = append(unionSingleMapVal[keysOneStr], cliSliceInfoSingle.val)
			unionCliInfo[i] = unionSingleMapVal
		}
	}
	// fmt.Println("unionInfo:", unionInfo, "singleUvals ================:\n", unionInfo[2], len(unionInfo))
	// for i, maoinf := range unionInfo[2] {
	// 	fmt.Println("i:", i, "maoinf:\n", maoinf)
	// }
	return unionSvrInfo, unionCliInfo
}

func makeSingleUkeys(unionKeys [][]string) [][]valInfo {
	singleUkeys := make([][]valInfo, len(unionKeys))
	for i, unionKey := range unionKeys {
		singleUkeys[i] = make([]valInfo, len(unionKey))
	}
	return singleUkeys
}

func makeSingleUkeysStr(unionKeys [][]string) [][]string {
	singleUkeys := make([][]string, len(unionKeys))
	for i, unionKey := range unionKeys {
		singleUkeys[i] = make([]string, len(unionKey))
	}
	return singleUkeys
}

func valInfo2Str(a valInfo) string {
	return fmt.Sprintf("%v_%v_%v", a.Name, a.Type, a.Val)
}

func keys2Str(strlist []string) string {
	if len(strlist) < 2 {
		return strlist[0]
	}
	restr := strlist[0]
	for _, str := range strlist[1:] {
		restr = fmt.Sprintf("%v|%v", restr, str)
	}
	return restr
}

func unionKeysMaxLen(unionKeys [][]string) int {
	max := 0
	for _, s := range unionKeys {
		if len(s) > max {
			max = len(s)
		}
	}
	return max
}

func inUnionKey(target string, strList []string) (bool, int) {
	found := false
	index := 0
	for i, str := range strList {
		if str == target {
			found = true
			index = i
			break
		}
	}
	return found, index
}

func splitKeys(str string) [][]string {
	stack := []rune{}
	resultKeys := [][]string{}
	for _, char := range str {
		if char == '[' {
			stack = append(stack, char)
		} else if char == ']' {
			temp := []rune{}
			for len(stack) > 0 && stack[len(stack)-1] != '[' {
				temp = append(temp, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
			if len(stack) > 0 && stack[len(stack)-1] == ',' {
				stack = stack[:len(stack)-1]
			}
			resultKey := spliteString(string(temp))
			if len(resultKey) > 0 {
				resultKeys = append(resultKeys, resultKey)
			}
		} else if char == ' ' {

		} else {
			stack = append(stack, char)
		}
	}

	return resultKeys
}

func spliteString(s string) []string {
	result := []string{}
	stack := []rune{}
	for _, char := range reverseString(s) {
		if char != ',' {
			stack = append(stack, char)
		} else {
			result = append(result, string(stack))
			stack = []rune{}
		}
	}
	if len(stack) > 0 {
		result = append(result, string(stack))
	}
	return result
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func Test1prop() {
	// 初始化声明一个 map
	myMap := map[string]int{
		"key1": 10,
		"key2": 20,
	}

	// 判断 map 中是否存在指定的键
	key := "key3"
	defaultValue := 0

	if value, ok := myMap[key]; ok {
		fmt.Printf("键 %s 的值为 %d\n", key, value)
	} else {
		fmt.Printf("键 %s 不存在，赋予默认值 %d\n", key, defaultValue)
		myMap[key] = defaultValue
	}

	// 打印更新后的 map
	fmt.Println(myMap)
}

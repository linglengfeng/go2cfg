package tocfg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tealeg/xlsx"
)

func Start(Options StartOptions) {
	Options.Default()
	err := createDir(Options.DirOutSvr)
	if err != nil {
		fmt.Println("CreateDir error:", err)
	}
	err = createDir(Options.DirOutCli)
	if err != nil {
		fmt.Println("CreateDir error:", err)
	}
	excelFiles, err := os.ReadDir(Options.DirExcel)
	if err != nil {
		fmt.Println("read dir error:", err)
	}
	for _, file := range excelFiles {
		if isTempFile(file.Name()) {
			continue
		}
		startSingle(file.Name(), Options)
	}
}

func startSingle(fileName string, Options StartOptions) {
	singleFileName, err := xlsx.OpenFile(Options.DirExcel + "/" + fileName)
	if err != nil {
		fmt.Println("open Excel:", singleFileName, "err:", err)
		return
	}
	singleFileBase := filepath.Base(fileName)
	excelFileName := strings.TrimSuffix(singleFileBase, filepath.Ext(singleFileBase))
	for _, sheet := range singleFileName.Sheets {
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
		// 设置 golbal prop
		setPropVal(excelFileName, sheet.Name, sheet.Rows)
		// 主键信息
		svrPrimarykeyInfo := []PrimarykeyVal{}
		cliPrimarykeyInfo := []PrimarykeyVal{}
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
			// 增加的单行数据
			addSvrValInfoList := ValInfoList{}
			addCliValInfoList := ValInfoList{}
			for cellIndex, cell := range row.Cells {
				text := cell.String()
				addkeyName := sheet.Rows[globalKey.Raw].Cells[cellIndex].String()
				addkeyType := sheet.Rows[globalType.Raw].Cells[cellIndex].String()
				if sheet.Rows[globalOutSvr.Raw].Cells[cellIndex].String() == TRUE { //是否导出
					addSvrValInfoList = append(addSvrValInfoList, ValInfo{Type: addkeyType, Val: text, Name: addkeyName})
				}
				if sheet.Rows[globalOutCli.Raw].Cells[cellIndex].String() == TRUE { //是否导出
					addCliValInfoList = append(addCliValInfoList, ValInfo{Type: addkeyType, Val: text, Name: addkeyName})
				}
			}
			idVal := row.Cells[VALUE_CELL+1].String()
			idType := sheet.Rows[globalType.Raw].Cells[VALUE_CELL+1].String()
			idName := sheet.Rows[globalName.Raw].Cells[VALUE_CELL+1].String()
			// 添加到主键信息中
			svrPrimarykeyInfo = append(svrPrimarykeyInfo, PrimarykeyVal{Key: ValInfo{Type: idType, Val: idVal, Name: idName}, ValList: addSvrValInfoList})
			cliPrimarykeyInfo = append(cliPrimarykeyInfo, PrimarykeyVal{Key: ValInfo{Type: idType, Val: idVal, Name: idName}, ValList: addCliValInfoList})
		}
		// 生成多键信息
		svrMapListInfo := generateUnionKeysInfo(svrPrimarykeyInfo)
		cliMapListInfo := generateUnionKeysInfo(cliPrimarykeyInfo)

		jsWriter := Writer{
			SvrFileName:       toString(globalExportSvr.Val),
			CliFileName:       toString(globalExportCli.Val),
			SvrPath:           Options.DirOutSvr,
			CliPath:           Options.DirOutCli,
			PrimarykeySvrData: svrPrimarykeyInfo,
			PrimarykeyCliData: cliPrimarykeyInfo,
			UnionKeysSvrData:  svrMapListInfo,
			UnionKeysCliData:  cliMapListInfo,
			Suffix:            ".json",
		}
		err = writeAll(jsWriter)
		if err != nil {
			fmt.Println("file:", excelFileName, "sheet name:", sheet.Name, "write file error:", err)
			return
		}

		fmt.Println("file:", excelFileName, "sheet name:", sheet.Name, "ok.")
	}
}

func writeAll(Writer Writer) error {
	// server primarykey data
	svrFile := Writer.SvrPath + "/" + Writer.SvrFileName + Writer.Suffix
	if !addCheckFileName(svrFile) {
		return fmt.Errorf("filename repeat:%v", svrFile)
	}
	svrdataPrimarykey := PrimarykeyJsonStr(Writer.PrimarykeySvrData)
	err := os.WriteFile(svrFile, append([]byte(svrdataPrimarykey), byte('\n')), 0644)
	if err != nil {
		return err
	}

	// server unionkey data
	svrUnionFile := Writer.SvrPath + "/" + Writer.SvrFileName + "_ukey" + Writer.Suffix
	if !addCheckFileName(svrUnionFile) {
		return fmt.Errorf("filename repeat:%v", svrUnionFile)
	}
	svrdataUnionkeys := UnionKeysJsonStr(Writer.UnionKeysSvrData)
	err = os.WriteFile(svrUnionFile, append([]byte(svrdataUnionkeys), byte('\n')), 0644)
	if err != nil {
		return err
	}

	// client primarykey data
	cliFile := Writer.CliPath + "/" + Writer.CliFileName + Writer.Suffix
	if !addCheckFileName(cliFile) {
		return fmt.Errorf("filename repeat:%v", cliFile)
	}
	globalCheckSameFile[cliFile] = struct{}{}
	clidataPrimarykey := PrimarykeyJsonStr(Writer.PrimarykeyCliData)
	err = os.WriteFile(svrFile, append([]byte(clidataPrimarykey), byte('\n')), 0644)
	if err != nil {
		return err
	}

	// client unionkey data
	cliUnionFile := Writer.CliPath + "/" + Writer.CliFileName + "_ukey" + Writer.Suffix
	if !addCheckFileName(cliUnionFile) {
		return fmt.Errorf("filename repeat:%v", cliUnionFile)
	}
	clidataUnionkeys := UnionKeysJsonStr(Writer.UnionKeysCliData)
	err = os.WriteFile(cliUnionFile, append([]byte(clidataUnionkeys), byte('\n')), 0644)
	if err != nil {
		return err
	}
	return nil
}

func UnionKeysJsonStr(svrMapListInfo []map[string]ValInfoList2) string {
	unionKeys := splitKeys(fmt.Sprintf("%v", globalUnionSvrKeys.Val))
	data := "{\n"
	for kindex, mapv := range svrMapListInfo {
		mapstr := "\t"

		kfield := unionKeys[kindex]
		kfieldstr := ""
		if len(kfield) > 1 {
			kfieldstr += "\""
			for _, kfields := range kfield {
				kfieldstr += kfields
			}
			kfieldstr += "\""
		} else {
			kfieldstr += "\"" + kfield[0] + "\""
		}

		singlemapStr := ""
		singlemapStr = kfieldstr + ":{\n"

		mapvstring := ""
		i := 0
		for mapvK, mapvV := range mapv {
			mapvstring += ""
			mapvsingleStr := ""
			mapvsingleStr += "\"" + mapvK + "\"" + ":[{\n"
			for mapvvi, mapvv := range mapvV {
				if mapvvi > 0 {
					mapvsingleStr += "}, {" + strValInfo(mapvv)
				} else {
					mapvsingleStr += strValInfo(mapvv)
				}
			}
			if i == len(mapv)-1 {
				mapvstring += mapvsingleStr + "}]\n"
			} else {
				mapvstring += mapvsingleStr + "}],\n"
			}
			i++
		}

		if kindex == len(svrMapListInfo)-1 {
			singlemapStr += mapvstring + "}\n\n\n"
		} else {
			singlemapStr += mapvstring + "},\n\n\n"
		}

		mapstr += singlemapStr

		mapstr += ""
		data += mapstr
	}
	data += "}"
	return data
}

func strValInfo(valInfoList []ValInfo) string {
	vstring := ""
	for iv, v := range valInfoList {
		vvalstr := ""
		switch v.Type {
		case STR:
			vvalstr = "\"" + v.Val + "\""
			// fmt.Println("type:", v.Type, v.Val, vvalstr)
		default:
			vvalstr = v.Val
		}
		if len(valInfoList)-1 == iv {
			vstring += "\t\t\"" + v.Name + "\"" + ": " + vvalstr + ""
		} else {
			vstring += "\t\t\"" + v.Name + "\"" + ": " + vvalstr + ",\n"
		}
	}
	vstring += "\n"
	// fmt.Println("vstring:", vstring)
	return vstring
}

func PrimarykeyJsonStr(svrPrimarykeyInfo []PrimarykeyVal) string {
	data := "{\n"
	for kvi, kv := range svrPrimarykeyInfo {
		kvstring := "\t\"" + kv.Key.Val + "\"" + ": {\n"
		vstring := ""

		for iv, v := range kv.ValList {
			vvalstr := ""
			switch v.Type {
			case STR:
				vvalstr = "\"" + v.Val + "\""
				// fmt.Println("type:", v.Type, v.Val, vvalstr)
			default:
				vvalstr = v.Val
			}
			if len(kv.ValList)-1 == iv {
				vstring += "\t\t\"" + v.Name + "\"" + ": " + vvalstr + ""
			} else {
				vstring += "\t\t\"" + v.Name + "\"" + ": " + vvalstr + ",\n"
			}
		}
		vstring += "\n"
		if len(svrPrimarykeyInfo)-1 == kvi {
			kvstring += vstring + "\t}\n"
		} else {
			kvstring += vstring + "\t},\n"
		}

		data += kvstring
	}
	data += "}"
	return data
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

func generateUnionKeysInfo(primarykeyInfo []PrimarykeyVal) []map[string]ValInfoList2 {
	unionKeys := splitKeys(toString(globalUnionSvrKeys.Val))
	unionInfo := make([]map[string]ValInfoList2, len(unionKeys))
	for i := range unionKeys {
		unionInfo[i] = make(map[string]ValInfoList2)
	}
	for _, svrPrimarykeyInfoSingle := range primarykeyInfo {
		singleUkeys := makeSingleUkeys(unionKeys)
		singleUkeysStr := makeSingleUkeysStr(unionKeys)
		for _, singleVal := range svrPrimarykeyInfoSingle.ValList {
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
			unionSingleMapVal := unionInfo[i]
			keysOneStr := keys2Str(slist)
			unionSingleMapVal[keysOneStr] = append(unionSingleMapVal[keysOneStr], svrPrimarykeyInfoSingle.ValList)
			unionInfo[i] = unionSingleMapVal
		}
	}
	return unionInfo
}

// func generateUnionKeysInfo(svrPrimarykeyInfo []PrimarykeyVal, cliPrimarykeyInfo []PrimarykeyVal) ([]map[string]ValInfoList2, []map[string]ValInfoList2) {
// 	unionKeys := splitKeys(fmt.Sprintf("%v", globalUnionSvrKeys.Val))

// 	unionSvrInfo := make([]map[string]ValInfoList2, len(unionKeys))
// 	for i := range unionKeys {
// 		unionSvrInfo[i] = make(map[string]ValInfoList2)
// 	}
// 	for _, svrPrimarykeyInfoSingle := range svrPrimarykeyInfo {
// 		singleUkeys := makeSingleUkeys(unionKeys)       // 五个key ，，val 都是 ssvrPrimarykeyInfoSingle.val
// 		singleUkeysStr := makeSingleUkeysStr(unionKeys) // 方面用map，key直接转成string
// 		for _, singleVal := range svrPrimarykeyInfoSingle.ValList {
// 			for unionKeyIndex, unionKey := range unionKeys {
// 				isin, inUnionpos := inUnionKey(singleVal.Name, unionKey)
// 				if isin {
// 					singleUkeys[unionKeyIndex][inUnionpos] = singleVal
// 					singleUkeysStr[unionKeyIndex][inUnionpos] = valInfo2Str(singleVal)

// 				}
// 			}
// 		}
// 		// keys 生成结束
// 		for i, slist := range singleUkeysStr {
// 			unionSingleMapVal := unionSvrInfo[i]
// 			keysOneStr := keys2Str(slist)
// 			unionSingleMapVal[keysOneStr] = append(unionSingleMapVal[keysOneStr], svrPrimarykeyInfoSingle.ValList)
// 			unionSvrInfo[i] = unionSingleMapVal
// 		}
// 	}

// 	// client
// 	unionCliInfo := make([]map[string]([][]ValInfo), len(unionKeys))
// 	for i := range unionKeys {
// 		unionCliInfo[i] = make(map[string]([][]ValInfo))
// 	}
// 	for _, cliPrimarykeyInfoSingle := range cliPrimarykeyInfo {
// 		singleUkeys := makeSingleUkeys(unionKeys)       // 五个key ，，val 都是 svrPrimarykeyInfoSingle.val
// 		singleUkeysStr := makeSingleUkeysStr(unionKeys) // 方面用map，key直接转成string
// 		for _, singleVal := range cliPrimarykeyInfoSingle.ValList {
// 			for unionKeyIndex, unionKey := range unionKeys {
// 				isin, inUnionpos := inUnionKey(singleVal.Name, unionKey)
// 				if isin {
// 					singleUkeys[unionKeyIndex][inUnionpos] = singleVal
// 					singleUkeysStr[unionKeyIndex][inUnionpos] = valInfo2Str(singleVal)

// 				}
// 			}
// 		}
// 		// keys 生成结束
// 		for i, slist := range singleUkeysStr {
// 			unionSingleMapVal := unionCliInfo[i]
// 			keysOneStr := keys2Str(slist)
// 			unionSingleMapVal[keysOneStr] = append(unionSingleMapVal[keysOneStr], cliPrimarykeyInfoSingle.ValList)
// 			unionCliInfo[i] = unionSingleMapVal
// 		}
// 	}
// 	// fmt.Println("unionInfo:", unionInfo, "singleUvals ================:\n", unionInfo[2], len(unionInfo))
// 	// for i, maoinf := range unionInfo[2] {
// 	// 	fmt.Println("i:", i, "maoinf:\n", maoinf)
// 	// }
// 	return unionSvrInfo, unionCliInfo
// }

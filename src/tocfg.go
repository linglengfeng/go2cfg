package tocfg

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/tealeg/xlsx"
)

func Start(Options StartOptions) {
	Options.GenDeal()
	if len(Options.WritersOpts) < 1 {
		fmt.Println("no writer,dont deal, finished.")
		return
	}
	excelFiles, err := os.ReadDir(Options.DirExcel)
	if err != nil {
		fmt.Println("read dir error:", err)
	}

	writerlist, wokerlist := preStartSingle(Options)
	if len(writerlist) <= 0 {
		return
	}
	// for _, file := range excelFiles {
	// 	if isTempFile(file.Name()) {
	// 		continue
	// 	}
	// 	startSingle(file.Name(), Options, wokerlist, writerlist)
	// }
	StartRecursion(len(excelFiles), excelFiles, Options, wokerlist, writerlist)
}

func startSingle(fileName string, Options StartOptions, wokerlist []WriterWorker, writerlist []*Writer) {
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
		svrMapListInfo := generateUnionKeysInfo(splitKeys(toString(globalUnionSvrKeys.Val)), svrPrimarykeyInfo)
		cliMapListInfo := generateUnionKeysInfo(splitKeys(toString(globalUnionCliKeys.Val)), cliPrimarykeyInfo)

		writeAll(excelFileName, sheet.Name, wokerlist, writerlist, svrPrimarykeyInfo, cliPrimarykeyInfo, svrMapListInfo, cliMapListInfo)
	}
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

func generateUnionKeysInfo(unionKeys [][]string, primarykeyInfo []PrimarykeyVal) []map[string]ValInfoList2 {
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

func StartRecursion(len int, excelFiles []fs.DirEntry, Options StartOptions, wokerlist []WriterWorker, writerlist []*Writer) {
	if len <= 0 {
		return
	}
	file := excelFiles[len-1]
	if !isTempFile(file.Name()) {
		startSingle(file.Name(), Options, wokerlist, writerlist)
	}
	StartRecursion(len-1, excelFiles, Options, wokerlist, writerlist)
}

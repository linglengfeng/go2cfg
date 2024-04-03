package tocfg

import (
	"fmt"
	"os"
	"strings"

	"github.com/tealeg/xlsx"
)

func createDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0644)
	}
	return nil
}

func isTempFile(name string) bool {
	if len(name) > 0 {
		return name[0] == '$' || name[0] == '~'
	}
	return false
}

func setPropVal(fileName string, sheetName string, rows []*xlsx.Row) {
	for _, prop := range globalCheckList {
		cells := rows[prop.Raw].Cells
		name := ""
		switch prop.Type {
		case EXPORT_SVR:
			if len(cells)-1 >= prop.Name {
				name = rows[prop.Raw].Cells[prop.Name].String()
			} else {
				name = fileName + "_" + sheetName
			}
			globalExportSvr.Val = name
		case EXPORT_CLI:
			if len(cells)-1 >= prop.Name {
				name = rows[prop.Raw].Cells[prop.Name].String()
			} else {
				name = fileName + "_" + sheetName
			}
			globalExportCli.Val = name
		case PRIMARY_KEY:
			globalPrimaryKey.Val = rows[prop.Raw].Cells[prop.Name].String()
		case UNION_KEYS_SVR:
			if len(cells)-1 >= prop.Name {
				name = rows[prop.Raw].Cells[prop.Name].String()
			}
			globalUnionSvrKeys.Val = name
		case UNION_KEYS_CLI:
			if len(cells)-1 >= prop.Name {
				name = rows[prop.Raw].Cells[prop.Name].String()
			}
			globalUnionCliKeys.Val = name
		case OUT_SVR:
			globalOutSvrKeys = ValInfoList{}
			if len(cells) > 1 {
				for icell, cell := range cells {
					if cell.String() == TRUE {
						globalOutSvrKeys = append(globalOutSvrKeys, ValInfo{Name: rows[globalKey.Raw].Cells[icell].String(), Type: rows[globalType.Raw].Cells[icell].String()})
					}
				}
			}
		case OUT_CLI:
			globalOutCliKeys = ValInfoList{}
			if len(cells) > 1 {
				for icell, cell := range cells {
					if cell.String() == TRUE {
						globalOutCliKeys = append(globalOutCliKeys, ValInfo{Name: rows[globalKey.Raw].Cells[icell].String(), Type: rows[globalType.Raw].Cells[icell].String()})
					}
				}
			}
		default:
			continue
		}
	}
}

func toString(s any) string {
	return fmt.Sprintf("%v", s)
}

func makeSingleUkeys(unionKeys [][]string) [][]ValInfo {
	singleUkeys := make([][]ValInfo, len(unionKeys))
	for i, unionKey := range unionKeys {
		singleUkeys[i] = make([]ValInfo, len(unionKey))
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

func valInfo2Str(v ValInfo) string {
	return v.Name + "_" + v.Type + "_" + v.Val
}

func str2valInfo(str string) ValInfo {
	splitStr := strings.Split(str, "_")
	valstr := ""
	for _, s := range splitStr[2:] {
		valstr += s
	}
	newA := ValInfo{Name: splitStr[0], Type: splitStr[1], Val: valstr}
	return newA
}

func keys2Str(strlist []string) string {
	if len(strlist) < 2 {
		return strlist[0]
	}
	restr := strlist[0]
	for _, str := range strlist[1:] {
		restr = restr + "|" + str
	}
	return restr
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

func toStartWriter(writersStr any) []WritersOpts {
	writers := writersStr.([]any)
	writerlist := []WritersOpts{}
	for _, writer := range writers {
		writerMap := writer.(map[string]any)
		outCliDir := writerMap["out_cli_dir"]
		outSvrDir := writerMap["out_svr_dir"]
		filetype := writerMap["type"]
		if outCliDir == nil || outSvrDir == nil || filetype == nil {
			continue
		}
		writerlist = append(writerlist, WritersOpts{Type: filetype.(string), DirOutSvr: outSvrDir.(string), DirOutCli: outCliDir.(string)})
	}
	return writerlist
}

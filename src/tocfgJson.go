package tocfg

import (
	"fmt"
	"os"
	"strings"
)

func writeJson(writer Writer) error {
	// 修改文件后缀名
	writer.Suffix = ".json"

	// 创建目录
	createDir(writer.SvrPath)
	createDir(writer.CliPath)

	// server primarykey data
	svrFile := writer.SvrPath + "/" + writer.SvrFileName + writer.Suffix
	svrdataPrimarykey := PrimarykeyJsonStr(writer.PrimarykeySvrData)
	err := os.WriteFile(svrFile, append([]byte(svrdataPrimarykey), byte('\n')), 0644)
	if err != nil {
		return err
	}

	// server unionkey data
	svrUnionFile := writer.SvrPath + "/" + writer.SvrFileName + "_ukey" + writer.Suffix
	svrdataUnionkeys := UnionKeysJsonStr(writer.UnionKeysSvrData)
	err = os.WriteFile(svrUnionFile, append([]byte(svrdataUnionkeys), byte('\n')), 0644)
	if err != nil {
		return err
	}

	// client primarykey data
	cliFile := writer.CliPath + "/" + writer.CliFileName + writer.Suffix
	clidataPrimarykey := PrimarykeyJsonStr(writer.PrimarykeyCliData)
	err = os.WriteFile(cliFile, append([]byte(clidataPrimarykey), byte('\n')), 0644)
	if err != nil {
		return err
	}

	// client unionkey data
	cliUnionFile := writer.CliPath + "/" + writer.CliFileName + "_ukey" + writer.Suffix
	clidataUnionkeys := UnionKeysJsonStr(writer.UnionKeysCliData)
	err = os.WriteFile(cliUnionFile, append([]byte(clidataUnionkeys), byte('\n')), 0644)
	if err != nil {
		return err
	}
	return nil
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

func strValInfo(valInfoList []ValInfo) string {
	vstring := ""
	for iv, v := range valInfoList {
		vvalstr := valInfo2JsStr(v)
		if len(valInfoList)-1 == iv {
			vstring += "\t\t\"" + v.Name + "\"" + ": " + vvalstr + ""
		} else {
			vstring += "\t\t\"" + v.Name + "\"" + ": " + vvalstr + ",\n"
		}
	}
	vstring += "\n"
	return vstring
}

func valInfo2JsStr(v ValInfo) string {
	switch v.Type {
	case STR:
		return "\"" + v.Val + "\""
	default:
		return v.Val
	}
}

func UnionKeysJsonStr(svrMapListInfo []map[string]ValInfoList2) string {
	unionKeys := splitKeys(fmt.Sprintf("%v", globalUnionSvrKeys.Val))
	data := "{\n"
	for kindex, mapv := range svrMapListInfo {
		mapstr := "\t"

		kfield := unionKeys[kindex]
		kfieldstr := ""
		if len(kfield) > 1 {
			kfieldstr += "\"" + kfield[0]
			for _, kfields := range kfield[1:] {
				kfieldstr += "_" + kfields
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
			mapvsingleStr := "\t\t"
			mapvK1 := keystr2jsstr(mapvK)
			mapvsingleStr += mapvK1 + ":[{\n"
			for mapvvi, mapvv := range mapvV {
				if mapvvi > 0 {
					mapvsingleStr += "\t\t}, {\n" + strValInfo(mapvv)
				} else {
					mapvsingleStr += strValInfo(mapvv)
				}
			}
			if i == len(mapv)-1 {
				mapvstring += mapvsingleStr + "\t\t}]\n"
			} else {
				mapvstring += mapvsingleStr + "\t\t}],\n"
			}
			i++
		}

		if kindex == len(svrMapListInfo)-1 {
			singlemapStr += mapvstring + "\t}\n\n"
		} else {
			singlemapStr += mapvstring + "\t},\n\n"
		}

		mapstr += singlemapStr

		mapstr += ""
		data += mapstr
	}
	data += "}"
	return data
}

func keystr2jsstr(str string) string {
	splitStr := strings.Split(str, "|")
	valstr := ""
	for i, s := range splitStr {
		val := str2valInfo(s)
		if i == len(splitStr)-1 {
			valstr += val.Val
		} else {
			valstr += val.Val + "_"
		}
	}
	return "\"" + valstr + "\""
}

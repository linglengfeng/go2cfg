package tocfg

import (
	"os"
	"strings"
)

// ---------------------------- 接口 start ---------------------------------
type jsonWorker struct {
	writer *Writer
}

func (worker *jsonWorker) AddSuffix() {
	worker.writer.Suffix = ".json"
}
func (worker *jsonWorker) ClearSvrDir() {
	svrErlPath := worker.writer.SvrPath
	deleteDir(svrErlPath)
}
func (worker *jsonWorker) ClearCliDir() {
	cliErlPath := worker.writer.CliPath
	deleteDir(cliErlPath)
}
func (worker *jsonWorker) CreateSvrDir() {
	svrErlPath := worker.writer.SvrPath
	createDir(svrErlPath)
}
func (worker *jsonWorker) CreateCliDir() {
	cliErlPath := worker.writer.CliPath
	createDir(cliErlPath)
}
func (worker *jsonWorker) WriteSvrIn() {
	unionkeySuffix := "_ukey"
	path := worker.writer.SvrPath
	filename := worker.writer.SvrFileName
	exportUkeys := worker.writer.SvrExportUkeys
	primarykeyData := worker.writer.PrimarykeySvrData
	unionKeysData := worker.writer.UnionKeysSvrData
	writeJsonIn(unionkeySuffix, path, filename, exportUkeys, primarykeyData, unionKeysData, *worker.writer)
}
func (worker *jsonWorker) WriteCliIn() {
	unionkeySuffix := "_ukey"
	path := worker.writer.CliPath
	filename := worker.writer.CliFileName
	exportUkeys := worker.writer.CliExportUkeys
	primarykeyData := worker.writer.PrimarykeyCliData
	unionKeysData := worker.writer.UnionKeysCliData
	writeJsonIn(unionkeySuffix, path, filename, exportUkeys, primarykeyData, unionKeysData, *worker.writer)
}

// ---------------------------- 接口 end ---------------------------------

func writeJsonIn(unionkeySuffix string, path string, filename string, exportUkeys [][]string, primarykeyData []PrimarykeyVal, unionKeysData []map[string]ValInfoList2, writer Writer) error {
	// server primarykey data
	file := path + "/" + filename + writer.Suffix
	dataPrimarykey := PrimarykeyJsonStr(primarykeyData)
	err := os.WriteFile(file, append([]byte(dataPrimarykey), byte('\n')), 0644)
	if err != nil {
		return err
	}
	// server unionkey data
	if len(exportUkeys) > 0 {
		unionFileName := path + "/" + filename + unionkeySuffix + writer.Suffix
		dataUnionkeys := UnionKeysJsonStr(exportUkeys, unionKeysData)
		err = os.WriteFile(unionFileName, append([]byte(dataUnionkeys), byte('\n')), 0644)
		if err != nil {
			return err
		}
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

func UnionKeysJsonStr(unionKeys [][]string, mapListInfo []map[string]ValInfoList2) string {
	data := "{\n"
	for kindex, mapv := range mapListInfo {
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

		if kindex == len(mapListInfo)-1 {
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

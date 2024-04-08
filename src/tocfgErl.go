package tocfg

import (
	"os"
	"strings"
)

// ---------------------------- 接口 start ---------------------------------
type erlWorker struct {
	writer *Writer
}

func (worker *erlWorker) AddSuffix() {
	worker.writer.Suffix = ".erl"
}
func (worker *erlWorker) ClearSvrDir() {
	svrErlPath := worker.writer.SvrPath
	deleteDir(svrErlPath)
}
func (worker *erlWorker) ClearCliDir() {
	cliErlPath := worker.writer.CliPath
	deleteDir(cliErlPath)
}
func (worker *erlWorker) CreateSvrDir() {
	svrErlPath := worker.writer.SvrPath
	createDir(svrErlPath)
}
func (worker *erlWorker) CreateCliDir() {
	cliErlPath := worker.writer.CliPath
	createDir(cliErlPath)
}
func (worker *erlWorker) WriteSvrIn() {
	perfix := "cfg_"
	unionkeySuffix := "_ukey"
	path := worker.writer.SvrPath
	filename := worker.writer.SvrFileName
	exportKeys := worker.writer.SvrExportKeys
	exportUkeys := worker.writer.SvrExportUkeys
	primarykeyData := worker.writer.PrimarykeySvrData
	unionKeysData := worker.writer.UnionKeysSvrData
	writeErlIn(perfix, unionkeySuffix, path, filename, exportKeys, exportUkeys, primarykeyData, unionKeysData, *worker.writer)
}
func (worker *erlWorker) WriteCliIn() {
	perfix := "cfg_"
	unionkeySuffix := "_ukey"
	path := worker.writer.CliPath
	filename := worker.writer.CliFileName
	exportKeys := worker.writer.CliExportKeys
	exportUkeys := worker.writer.CliExportUkeys
	primarykeyData := worker.writer.PrimarykeyCliData
	unionKeysData := worker.writer.UnionKeysCliData
	writeErlIn(perfix, unionkeySuffix, path, filename, exportKeys, exportUkeys, primarykeyData, unionKeysData, *worker.writer)
}

// ---------------------------- 接口 end ---------------------------------

func writeErlIn(perfix string, unionkeySuffix string, path string, filename string, exportKeys ValInfoList,
	exportUkeys [][]string, primarykeyData []PrimarykeyVal, unionKeysData []map[string]ValInfoList2, writer Writer) error {
	// server primarykey data
	svrFile := path + "/" + perfix + filename + writer.Suffix
	modulename := perfix + filename
	svrdataPrimarykey := PrimarykeyErlStr(modulename, exportKeys, primarykeyData)
	err := os.WriteFile(svrFile, append([]byte(svrdataPrimarykey), byte('\n')), 0644)
	if err != nil {
		return err
	}
	// server unionkey data
	if len(exportUkeys) > 0 {
		svrUnionFile := path + "/" + perfix + filename + unionkeySuffix + writer.Suffix
		ukeyModuleName := modulename + unionkeySuffix
		svrdataUnionkeys := UnionKeysErlStr(ukeyModuleName, exportKeys, exportUkeys, unionKeysData)
		err = os.WriteFile(svrUnionFile, append([]byte(svrdataUnionkeys), byte('\n')), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func PrimarykeyErlStr(modulename string, keys ValInfoList, svrPrimarykeyInfo []PrimarykeyVal) string {
	// module
	module := "-module(" + modulename + ").\n\n"

	// record
	recordname := "cfg"
	record := "-record(" + recordname + ", {"
	for ivalinfo, valinfo := range keys {
		if ivalinfo == len(keys)-1 {
			record += valinfo.Name
		} else {
			record += valinfo.Name + ", "
		}
	}
	record += "}).\n\n"

	// func
	funcname := "get"
	funcstr := ""
	for _, kvinfo := range svrPrimarykeyInfo {
		kstr := valInfo2ErlStr(kvinfo.Key)
		vstr := valInfoList2ErlRecordStr(recordname, kvinfo.ValList)
		funcstr += funcname + "(" + kstr + ") ->\n\t" + vstr + ";\n"
	}
	funcstr += erlUndefined(funcname)

	// get
	getstr := ""
	exportfuncstr := ""
	for ivalinfo, valinfo := range keys {
		getfuncname := funcname + "_" + valinfo.Name
		getstr += getfuncname + "(Val=#" + recordname + "{}) ->\n" + "\tVal#" + recordname + "." + valinfo.Name + ";\n" + erlUndefined(getfuncname)
		if ivalinfo == len(keys)-1 {
			exportfuncstr += getfuncname + "/1"
		} else {
			exportfuncstr += getfuncname + "/1, "
		}

	}

	// export
	export := "-export[" + funcname + "/1, " + exportfuncstr + "].\n\n"

	return module + export + record + funcstr + getstr
}

func UnionKeysErlStr(ukeyModuleName string, keys ValInfoList, exportUkeys [][]string, unionKeysSvrData []map[string]ValInfoList2) string {
	// module
	module := "-module(" + ukeyModuleName + ").\n\n"

	// record
	recordname := "cfg"
	record := "-record(" + recordname + ", {"
	for ivalinfo, valinfo := range keys {
		if ivalinfo == len(keys)-1 {
			record += valinfo.Name
		} else {
			record += valinfo.Name + ", "
		}
	}
	record += "}).\n\n"

	// func
	funcstr := ""
	funcnamelist := ""
	for uki, ukv := range exportUkeys {
		funcname := "get_"
		for ukvstri, ukvstr := range ukv {
			if len(ukv)-1 == ukvstri {
				funcname += ukvstr
			} else {
				funcname += ukvstr + "_"
			}
		}
		if uki == len(exportUkeys)-1 {
			funcnamelist += funcname + "/1"
		} else {
			funcnamelist += funcname + "/1, "
		}

		singleMapInfo := unionKeysSvrData[uki]
		for singleMapInfok, singleMapInfov := range singleMapInfo {
			vstr := "["
			for singleMapInfovvi, singleMapInfovv := range singleMapInfov {
				if singleMapInfovvi == len(singleMapInfov)-1 {
					vstr += valInfoList2ErlRecordStr(recordname, singleMapInfovv)
				} else {
					vstr += valInfoList2ErlRecordStr(recordname, singleMapInfovv) + ",\n\t"
				}
			}
			vstr += "]"
			singleMapInfok1 := keystr2Erlstr(singleMapInfok)
			funcstr += funcname + "({" + singleMapInfok1 + "}) ->\n\t" + vstr + ";\n"
		}
		funcstr += erlUndefinedDef(funcname, "[]")

	}

	// export
	export := "-export[" + funcnamelist + "].\n\n"

	return module + export + record + funcstr
}

func valInfo2ErlStr(v ValInfo) string {
	switch v.Type {
	case INT:
		return v.Val
	case STR:
		return "\"" + v.Val + "\""
	case LIST:
		return v.Val
	default:
		return v.Val
	}
}

func valInfoList2ErlRecordStr(recordname string, vlist ValInfoList) string {
	record := "#" + recordname + "{"
	for i, v := range vlist {

		vstr := valInfo2ErlStr(v)
		if i == len(vlist)-1 {
			record += v.Name + "=" + vstr + ""
		} else {
			record += v.Name + "=" + vstr + ","
		}
	}
	record += "}"
	return record
}

func keystr2Erlstr(str string) string {
	splitStr := strings.Split(str, "|")
	valstr := ""
	for i, s := range splitStr {
		val := str2valInfo(s)
		if i == len(splitStr)-1 {
			valstr += valInfo2ErlStr(val)
		} else {
			valstr += valInfo2ErlStr(val) + ", "
		}
	}
	return valstr
}

func erlUndefined(funcname string) string {
	return funcname + "(_) ->\n\tundefined.\n\n"
}

func erlUndefinedDef(funcname string, defstr string) string {
	return funcname + "(_) ->\n\t" + defstr + ".\n\n"
}

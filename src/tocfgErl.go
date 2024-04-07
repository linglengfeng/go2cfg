package tocfg

import (
	"os"
	"strings"
)

func writeErl(writer Writer) error {
	// 修改文件后缀名
	writer.Suffix = ".erl"
	perfix := "cfg_"
	unionkeySuffix := "_ukey"

	if writer.SvrPath != "" {
		// 创建目录
		createDir(writer.SvrPath)
		// server primarykey data
		svrFile := writer.SvrPath + "/" + perfix + writer.SvrFileName + writer.Suffix
		modulename := perfix + writer.SvrFileName
		svrdataPrimarykey := PrimarykeyErlStr(modulename, writer.SvrExportKeys, writer.PrimarykeySvrData)
		err := os.WriteFile(svrFile, append([]byte(svrdataPrimarykey), byte('\n')), 0644)
		if err != nil {
			return err
		}
		// server unionkey data
		if len(writer.SvrExportUkeys) > 0 {
			svrUnionFile := writer.SvrPath + "/" + perfix + writer.SvrFileName + unionkeySuffix + writer.Suffix
			ukeyModuleName := modulename + unionkeySuffix
			svrdataUnionkeys := UnionKeysErlStr(ukeyModuleName, writer.SvrExportKeys, writer.SvrExportUkeys, writer.UnionKeysSvrData)
			err = os.WriteFile(svrUnionFile, append([]byte(svrdataUnionkeys), byte('\n')), 0644)
			if err != nil {
				return err
			}
		}
	}

	if writer.CliPath != "" {
		// 创建目录
		createDir(writer.CliPath)
		// client primarykey data
		cliFile := writer.CliPath + "/" + perfix + writer.CliFileName + writer.Suffix
		modulename := perfix + writer.CliFileName
		clidataPrimarykey := PrimarykeyErlStr(modulename, writer.SvrExportKeys, writer.PrimarykeySvrData)
		err := os.WriteFile(cliFile, append([]byte(clidataPrimarykey), byte('\n')), 0644)
		if err != nil {
			return err
		}
		// client unionkey data
		if len(writer.CliExportUkeys) > 0 {
			cliUnionFile := writer.CliPath + "/" + perfix + writer.CliFileName + unionkeySuffix + writer.Suffix
			ukeyModuleName := modulename + unionkeySuffix
			clidataUnionkeys := UnionKeysErlStr(ukeyModuleName, writer.CliExportKeys, writer.CliExportUkeys, writer.UnionKeysCliData)
			err = os.WriteFile(cliUnionFile, append([]byte(clidataUnionkeys), byte('\n')), 0644)
			if err != nil {
				return err
			}
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

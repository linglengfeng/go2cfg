package tocfg

import (
	"os"
)

func writeErl(writer Writer) error {
	// 修改文件后缀名
	writer.Suffix = ".erl"

	// 创建目录
	createDir(writer.SvrPath)
	createDir(writer.CliPath)

	perfix := "cfg_"
	// unionkeySuffix := "_ukey"

	// server primarykey data
	svrFile := writer.SvrPath + "/" + perfix + writer.SvrFileName + writer.Suffix
	modulename := perfix + writer.SvrFileName
	svrdataPrimarykey := PrimarykeyErlStr(modulename, writer.SvrExportKeys, writer.PrimarykeySvrData)
	err := os.WriteFile(svrFile, append([]byte(svrdataPrimarykey), byte('\n')), 0644)
	if err != nil {
		return err
	}

	// // server unionkey data
	// svrUnionFile :=  writer.SvrPath + "/" + perfix + writer.SvrFileName + unionkeySuffix + writer.Suffix
	// svrdataUnionkeys := UnionKeysJsonStr(writer.UnionKeysSvrData)
	// err = os.WriteFile(svrUnionFile, append([]byte(svrdataUnionkeys), byte('\n')), 0644)
	// if err != nil {
	// 	return err
	// }

	// // client primarykey data
	// cliFile := writer.CliPath + "/" + perfix + writer.CliFileName + writer.Suffix
	// clidataPrimarykey := PrimarykeyJsonStr(writer.PrimarykeyCliData)
	// err = os.WriteFile(cliFile, append([]byte(clidataPrimarykey), byte('\n')), 0644)
	// if err != nil {
	// 	return err
	// }

	// // client unionkey data
	// cliUnionFile :=  writer.CliPath + "/" + perfix + writer.CliFileName + unionkeySuffix + writer.Suffix
	// clidataUnionkeys := UnionKeysJsonStr(writer.UnionKeysCliData)
	// err = os.WriteFile(cliUnionFile, append([]byte(clidataUnionkeys), byte('\n')), 0644)
	// if err != nil {
	// 	return err
	// }
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
	for ikvinfo, kvinfo := range svrPrimarykeyInfo {
		kstr := valInfo2ErlStr(kvinfo.Key)
		vstr := valInfoList2ErlRecordStr(recordname, kvinfo.ValList)
		if ikvinfo == len(svrPrimarykeyInfo)-1 {
			funcstr += funcname + "(" + kstr + ") ->\n\t" + vstr + ".\n\n"
		} else {
			funcstr += funcname + "(" + kstr + ") ->\n\t" + vstr + ";\n"
		}
	}

	// get
	getstr := ""
	exportfuncstr := ""
	for ivalinfo, valinfo := range keys {
		getfuncname := funcname + "_" + valinfo.Name
		getstr += getfuncname + "(Val=#" + recordname + "{}) ->\n" + "\tVal#" + recordname + "." + valinfo.Name + ".\n\n"
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
			record += v.Name + "=" + vstr + ", "
		}
	}
	record += "}"
	return record
}

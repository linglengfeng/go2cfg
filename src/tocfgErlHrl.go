package tocfg

import (
	"os"
)

func writeErlHrl(writer Writer) error {
	// 修改文件后缀名
	writer.Suffix = ".erl"
	perfix := "cfg_"
	unionkeySuffix := "_ukey"
	svrErlPath := writer.SvrPath + "/erl"
	svrHrlPath := writer.SvrPath + "/include"
	svrHrlFileName := svrHrlPath + "/cfg.hrl"
	svrErlInclude := "../include/cfg.hrl"
	if writer.SvrPath != "" {
		// 创建目录
		createDir(svrErlPath)
		createDir(svrHrlPath)
		// primarykey data
		svrFile := svrErlPath + "/" + perfix + writer.SvrFileName + writer.Suffix
		modulename := perfix + writer.SvrFileName
		svrdataPrimarykey, svrdataPrimarykeyRecord := PrimarykeyErlHrlStr(svrErlInclude, modulename, writer.SvrExportKeys, writer.PrimarykeySvrData)
		err := os.WriteFile(svrFile, append([]byte(svrdataPrimarykey), byte('\n')), 0644)
		if err != nil {
			return err
		}
		// unionkey data
		svrUkeysRecord := ""
		if len(writer.SvrExportUkeys) > 0 {
			svrUnionFile := svrErlPath + "/" + perfix + writer.SvrFileName + unionkeySuffix + writer.Suffix
			ukeyModuleName := modulename + unionkeySuffix
			svrdataUnionkeys, svrdataUnionkeysRecord := UnionKeysErlHrlStr(svrErlInclude, ukeyModuleName, writer.SvrExportKeys, writer.SvrExportUkeys, writer.UnionKeysSvrData)
			svrUkeysRecord = svrdataUnionkeysRecord
			err = os.WriteFile(svrUnionFile, append([]byte(svrdataUnionkeys), byte('\n')), 0644)
			if err != nil {
				return err
			}
		}
		// hrl data
		erlfilehrl, err := os.OpenFile(svrHrlFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer erlfilehrl.Close()
		_, err = erlfilehrl.WriteString(svrdataPrimarykeyRecord + svrUkeysRecord)
		if err != nil {
			return err
		}
	}

	cliErlPath := writer.CliPath + "/erl"
	cliHrlPath := writer.CliPath + "/include"
	cliHrlFileName := cliHrlPath + "/cfg.hrl"
	cliErlInclude := "../include/cfg.hrl"
	if writer.CliPath != "" {
		// 创建目录
		createDir(cliErlPath)
		createDir(cliHrlPath)
		// primarykey data
		cliFile := cliErlPath + "/" + perfix + writer.CliFileName + writer.Suffix
		modulename := perfix + writer.CliFileName
		clidataPrimarykey, clidataPrimarykeyRecord := PrimarykeyErlHrlStr(cliErlInclude, modulename, writer.CliExportKeys, writer.PrimarykeyCliData)
		err := os.WriteFile(cliFile, append([]byte(clidataPrimarykey), byte('\n')), 0644)
		if err != nil {
			return err
		}
		// unionkey data
		cliUkeysRecord := ""
		if len(writer.CliExportUkeys) > 0 {
			cliUnionFile := cliErlPath + "/" + perfix + writer.CliFileName + unionkeySuffix + writer.Suffix
			ukeyModuleName := modulename + unionkeySuffix
			clidataUnionkeys, clidataUnionkeysRecord := UnionKeysErlHrlStr(cliErlInclude, ukeyModuleName, writer.CliExportKeys, writer.CliExportUkeys, writer.UnionKeysCliData)
			cliUkeysRecord = clidataUnionkeysRecord
			err = os.WriteFile(cliUnionFile, append([]byte(clidataUnionkeys), byte('\n')), 0644)
			if err != nil {
				return err
			}
		}
		// hrl data
		erlfilehrl, err := os.OpenFile(cliHrlFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer erlfilehrl.Close()
		_, err = erlfilehrl.WriteString(clidataPrimarykeyRecord + cliUkeysRecord)
		if err != nil {
			return err
		}
	}
	return nil
}

func PrimarykeyErlHrlStr(erlInclude string, modulename string, keys ValInfoList, primarykeyInfo []PrimarykeyVal) (string, string) {
	// module
	module := "-module(" + modulename + ").\n\n"

	// record
	recordname := modulename
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
	for _, kvinfo := range primarykeyInfo {
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

	// include
	include := "-include(\"" + erlInclude + "\").\n\n"

	return (module + export + include + funcstr + getstr), record
}

func UnionKeysErlHrlStr(erlInclude string, ukeyModuleName string, keys ValInfoList, exportUkeys [][]string, unionKeysData []map[string]ValInfoList2) (string, string) {
	// module
	module := "-module(" + ukeyModuleName + ").\n\n"

	// record
	recordname := ukeyModuleName
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

		singleMapInfo := unionKeysData[uki]
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

	// include
	include := "-include(\"" + erlInclude + "\").\n\n"

	return (module + export + include + funcstr), record
}

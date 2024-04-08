package tocfg

import (
	"os"
)

// ---------------------------- 接口 start ---------------------------------
type erlHrlWorker struct {
	writer     *Writer
	svrErlPath string
	svrHrlPath string
	cliErlPath string
	cliHrlPath string
}

func (worker *erlHrlWorker) AddSuffix() {
	worker.writer.Suffix = ".erl"
}
func (worker *erlHrlWorker) ClearSvrDir() {
	svrErlPath := worker.writer.SvrPath + "/erl"
	svrHrlPath := worker.writer.SvrPath + "/include"
	worker.svrErlPath = svrErlPath
	worker.svrHrlPath = svrHrlPath
	deleteDir(svrErlPath)
	deleteDir(svrHrlPath)
}
func (worker *erlHrlWorker) ClearCliDir() {
	cliErlPath := worker.writer.CliPath + "/erl"
	cliHrlPath := worker.writer.CliPath + "/include"
	worker.cliErlPath = cliErlPath
	worker.cliHrlPath = cliHrlPath
	deleteDir(cliErlPath)
	deleteDir(cliHrlPath)
}
func (worker *erlHrlWorker) CreateSvrDir() {
	svrErlPath := worker.svrErlPath
	svrHrlPath := worker.svrHrlPath
	createDir(svrErlPath)
	createDir(svrHrlPath)
}
func (worker *erlHrlWorker) CreateCliDir() {
	cliErlPath := worker.cliErlPath
	cliHrlPath := worker.cliHrlPath
	createDir(cliErlPath)
	createDir(cliHrlPath)
}
func (worker *erlHrlWorker) WriteSvrIn() {
	writeSvrInErlHrl(*worker)
}
func (worker *erlHrlWorker) WriteCliIn() {
	writeCliInErlHrl(*worker)
}

// ---------------------------- 接口 end ---------------------------------

func writeSvrInErlHrl(worker erlHrlWorker) {
	perfix := "cfg_"
	unionkeySuffix := "_ukey"
	erlInclude := "../include/cfg.hrl"
	hrlFileName := worker.svrHrlPath + "/cfg.hrl"
	erlPath := worker.svrErlPath
	writeErlHrlFile(perfix, unionkeySuffix, erlPath, hrlFileName, erlInclude, worker.writer)
}
func writeCliInErlHrl(worker erlHrlWorker) {
	perfix := "cfg_"
	unionkeySuffix := "_ukey"
	erlInclude := "../include/cfg.hrl"
	hrlFileName := worker.cliHrlPath + "/cfg.hrl"
	erlPath := worker.cliErlPath
	writeErlHrlFile(perfix, unionkeySuffix, erlPath, hrlFileName, erlInclude, worker.writer)
}

func writeErlHrlFile(perfix string, unionkeySuffix string, erlPath string, hrlFileName string, erlInclude string, writer *Writer) error {
	// primarykey data
	cliFile := erlPath + "/" + perfix + writer.CliFileName + writer.Suffix
	modulename := perfix + writer.CliFileName
	clidataPrimarykey, clidataPrimarykeyRecord := PrimarykeyErlHrlStr(erlInclude, modulename, writer.CliExportKeys, writer.PrimarykeyCliData)
	err := os.WriteFile(cliFile, append([]byte(clidataPrimarykey), byte('\n')), 0644)
	if err != nil {
		return err
	}
	// unionkey data
	cliUkeysRecord := ""
	if len(writer.CliExportUkeys) > 0 {
		cliUnionFile := erlPath + "/" + perfix + writer.CliFileName + unionkeySuffix + writer.Suffix
		ukeyModuleName := modulename + unionkeySuffix
		clidataUnionkeys, clidataUnionkeysRecord := UnionKeysErlHrlStr(erlInclude, ukeyModuleName, writer.CliExportKeys, writer.CliExportUkeys, writer.UnionKeysCliData)
		cliUkeysRecord = clidataUnionkeysRecord
		err = os.WriteFile(cliUnionFile, append([]byte(clidataUnionkeys), byte('\n')), 0644)
		if err != nil {
			return err
		}
	}
	// hrl data
	erlfilehrl, err := os.OpenFile(hrlFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer erlfilehrl.Close()
	_, err = erlfilehrl.WriteString(clidataPrimarykeyRecord + cliUkeysRecord)
	if err != nil {
		return err
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

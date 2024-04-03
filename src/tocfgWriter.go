package tocfg

import "fmt"

func writeAll(excelFileName string, sheetName string, Options StartOptions, svrPrimarykeyInfo []PrimarykeyVal, cliPrimarykeyInfo []PrimarykeyVal, svrMapListInfo []map[string]ValInfoList2, cliMapListInfo []map[string]ValInfoList2) {
	for _, writerOpt := range Options.WritersOpts {
		writer := Writer{
			writerType:        writerOpt.Type,
			SvrExportKeys:     globalOutSvrKeys,
			CliExportKeys:     globalOutCliKeys,
			SvrFileName:       toString(globalExportSvr.Val),
			CliFileName:       toString(globalExportCli.Val),
			SvrPath:           writerOpt.DirOutSvr,
			CliPath:           writerOpt.DirOutCli,
			PrimarykeySvrData: svrPrimarykeyInfo,
			PrimarykeyCliData: cliPrimarykeyInfo,
			UnionKeysSvrData:  svrMapListInfo,
			UnionKeysCliData:  cliMapListInfo,
			Suffix:            ".txt",
		}
		switch writerOpt.Type {
		case EXPORT_TYPE_JSON:
			writeJson(writer)
		case EXPORT_TYPE_TXT:
			writeTxt(writer)
		case EXPORT_TYPE_ERL:
			writeErl(writer)
		default:
			fmt.Println("not case this file type, error type:", writerOpt.Type, " excelFileName:", excelFileName, "sheetName:", sheetName)
			continue
		}
		fmt.Println("write:", writer.writerType, ",file:", excelFileName, ",sheetName:", sheetName, ", write ok.")
	}
}

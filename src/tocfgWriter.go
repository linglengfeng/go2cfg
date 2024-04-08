package tocfg

import "fmt"

func initWorker(wer *Writer) (WriterWorker, error) {
	switch wer.WriterType {
	case EXPORT_TYPE_JSON:
		return &jsonWorker{writer: wer}, nil
	case EXPORT_TYPE_TXT:
		return &txtWorker{writer: wer}, nil
	case EXPORT_TYPE_ERL:
		return &erlWorker{writer: wer}, nil
	case EXPORT_TYPE_ERL_HRL:
		return &erlHrlWorker{writer: wer}, nil
	default:
		return &erlHrlWorker{}, fmt.Errorf("not type:%v", wer.WriterType)
	}
}

func writeAll(excelFileName string, sheetName string, wokerlist []WriterWorker, writerlist []*Writer, svrPrimarykeyInfo []PrimarykeyVal, cliPrimarykeyInfo []PrimarykeyVal, svrMapListInfo []map[string]ValInfoList2, cliMapListInfo []map[string]ValInfoList2) {
	for writeri, writer := range writerlist {
		writer.SvrExportKeys = globalOutSvrKeys
		writer.CliExportKeys = globalOutCliKeys
		writer.SvrExportUkeys = splitKeys(fmt.Sprintf("%v", globalUnionSvrKeys.Val))
		writer.CliExportUkeys = splitKeys(fmt.Sprintf("%v", globalUnionCliKeys.Val))
		writer.SvrFileName = toString(globalExportSvr.Val)
		writer.CliFileName = toString(globalExportCli.Val)
		writer.PrimarykeySvrData = svrPrimarykeyInfo
		writer.PrimarykeyCliData = cliPrimarykeyInfo
		writer.UnionKeysSvrData = svrMapListInfo
		writer.UnionKeysCliData = cliMapListInfo
		writerWorking := WriteWorking{}

		if writer.SvrPath != "" {
			writerWorking.workCreateSvrAndWrite(wokerlist[writeri])
		}
		if writer.CliPath != "" {
			writerWorking.workCreateCliAndWrite(wokerlist[writeri])
		}

		fmt.Println("write:", writer.WriterType, ",file:", excelFileName, ",sheetName:", sheetName, ", write ok.")
	}
}

func preStartSingle(Options StartOptions) ([]*Writer, []WriterWorker) {
	var writerlist []*Writer
	var wokerlist []WriterWorker
	for _, writerOpt := range Options.WritersOpts {
		singleWriter := Writer{
			WriterType: writerOpt.Type,
			SvrPath:    writerOpt.DirOutSvr,
			CliPath:    writerOpt.DirOutCli,
		}
		worker, err := initWorker(&singleWriter)
		if err == nil {
			wokerlist = append(wokerlist, worker)
			writerlist = append(writerlist, &singleWriter)
			writerWorking := WriteWorking{}
			if singleWriter.SvrPath != "" && singleWriter.CliPath == "" {
				continue
			}
			if singleWriter.SvrPath != "" {
				writerWorking.preWorkClearSvrDir(worker)
			}
			if singleWriter.CliPath != "" {
				writerWorking.preWorkClearCliDir(worker)
			}
		}
	}
	return writerlist, wokerlist
}

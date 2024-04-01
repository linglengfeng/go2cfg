package tocfg

import "encoding/json"

type prop struct {
	Type string
	Raw  int
	Cell int
	Is   int
	Name int
	Val  any
}

type ValInfo struct {
	Name string
	Type string
	Val  string
}

type ValInfoList []ValInfo

type ValInfoList2 [][]ValInfo

type PrimarykeyVal struct {
	Key     ValInfo
	ValList ValInfoList
}

type UnionKeysVals struct {
	Keys ValInfoList
	Vals ValInfoList2
}

// func start options
type StartOptions struct {
	DirExcel  string
	DirOutSvr string
	DirOutCli string
}

type Writer struct {
	SvrFileName       string
	CliFileName       string
	SvrPath           string
	CliPath           string
	Suffix            string
	PrimarykeySvrData []PrimarykeyVal
	PrimarykeyCliData []PrimarykeyVal
	UnionKeysSvrData  []map[string]ValInfoList2
	UnionKeysCliData  []map[string]ValInfoList2
}

func (StartOpts *StartOptions) Default() StartOptions {
	if StartOpts.DirExcel == "" {
		StartOpts.DirExcel = DEFALUT_EXCELDIR
	}
	if StartOpts.DirOutSvr == "" {
		StartOpts.DirOutSvr = DEFALUT_SVRDIR
	}
	if StartOpts.DirOutCli == "" {
		StartOpts.DirOutCli = DEFAULT_CLIDIR
	}
	return *StartOpts
}

// func (Writer *Writer) SvrName() string {
// 	return Writer.SvrPath + "/" + Writer.SvrNameStr + "." + Writer.SuffixStr
// }

// func (Writer *Writer) CliName() string {
// 	return Writer.CliPath + "/" + Writer.SvrNameStr + "." + Writer.SuffixStr
// }

// func (Writer *Writer) ToSvrData() string {
// 	return ToString(Writer.SvrData)
// }

// func (Writer *Writer) ToCliData() string {
// 	return ToString(Writer.SvrData)
// }

func ToString(param map[string]any) string {
	byteData, err := json.Marshal(param)
	if err != nil {
		panic(err)
	}
	dataString := string(byteData)
	return dataString
}

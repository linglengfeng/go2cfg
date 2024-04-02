package tocfg

import (
	"encoding/json"
)

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
	DirExcel       string
	WritersOptsStr any
	WritersOpts    []WritersOpts
}

type WritersOpts struct {
	Type      string
	DirOutSvr string
	DirOutCli string
}

type Writer struct {
	writerType        string
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

func (StartOpts *StartOptions) GenDeal() StartOptions {
	if StartOpts.DirExcel == "" {
		StartOpts.DirExcel = DEFALUT_EXCELDIR
	}
	if StartOpts.WritersOptsStr == nil {
		StartOpts.WritersOpts = []WritersOpts{}
	} else {
		StartOpts.WritersOpts = toStartWriter(StartOpts.WritersOptsStr)
	}
	return *StartOpts
}

func ToString(param map[string]any) string {
	byteData, err := json.Marshal(param)
	if err != nil {
		panic(err)
	}
	dataString := string(byteData)
	return dataString
}

package tocfg

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

type WritersOpts struct {
	Type      string
	DirOutSvr string
	DirOutCli string
}

type Writer struct {
	WriterType        string
	SvrExportKeys     ValInfoList
	CliExportKeys     ValInfoList
	SvrExportUkeys    [][]string
	CliExportUkeys    [][]string
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

type WriteWorking struct {
}

func (w WriteWorking) preWorkClearSvrDir(worker WriterWorker) {
	worker.ClearSvrDir()
}
func (w WriteWorking) preWorkClearCliDir(worker WriterWorker) {
	worker.ClearCliDir()
}
func (w WriteWorking) workCreateCliAndWrite(worker WriterWorker) {
	worker.AddSuffix()
	worker.CreateCliDir()
	worker.WriteCliIn()
}
func (w WriteWorking) workCreateSvrAndWrite(worker WriterWorker) {
	worker.AddSuffix()
	worker.CreateSvrDir()
	worker.WriteSvrIn()
}

type WriterWorker interface {
	AddSuffix()
	ClearSvrDir()
	ClearCliDir()
	CreateSvrDir()
	CreateCliDir()
	WriteSvrIn()
	WriteCliIn()
}

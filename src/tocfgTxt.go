package tocfg

// ---------------------------- 接口 start ---------------------------------
type txtWorker struct {
	writer *Writer
}

func (worker *txtWorker) AddSuffix() {
	worker.writer.Suffix = ".txt"
}
func (worker *txtWorker) ClearSvrDir() {
	svrErlPath := worker.writer.SvrPath
	deleteDir(svrErlPath)
}
func (worker *txtWorker) ClearCliDir() {
	cliErlPath := worker.writer.CliPath
	deleteDir(cliErlPath)
}
func (worker *txtWorker) CreateSvrDir() {
	svrErlPath := worker.writer.SvrPath
	createDir(svrErlPath)
}
func (worker *txtWorker) CreateCliDir() {
	cliErlPath := worker.writer.CliPath
	createDir(cliErlPath)
}
func (worker *txtWorker) WriteSvrIn() {
	unionkeySuffix := "_ukey"
	path := worker.writer.SvrPath
	filename := worker.writer.SvrFileName
	exportUkeys := worker.writer.SvrExportUkeys
	primarykeyData := worker.writer.PrimarykeySvrData
	unionKeysData := worker.writer.UnionKeysSvrData
	writeJsonIn(unionkeySuffix, path, filename, exportUkeys, primarykeyData, unionKeysData, *worker.writer)
}
func (worker *txtWorker) WriteCliIn() {
	unionkeySuffix := "_ukey"
	path := worker.writer.CliPath
	filename := worker.writer.CliFileName
	exportUkeys := worker.writer.CliExportUkeys
	primarykeyData := worker.writer.PrimarykeyCliData
	unionKeysData := worker.writer.UnionKeysCliData
	writeJsonIn(unionkeySuffix, path, filename, exportUkeys, primarykeyData, unionKeysData, *worker.writer)
}

// ---------------------------- 接口 end ---------------------------------

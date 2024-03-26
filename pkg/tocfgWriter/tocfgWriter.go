package tocfgWriter

import "encoding/json"

type Writer struct {
	SvrNameStr string
	CliNameStr string
	SvrPath    string
	CliPath    string
	SuffixStr  string
	SvrData    map[string]any
	CliData    map[string]any
}

func (Writer *Writer) SvrName() string {
	return Writer.SvrPath + "/" + Writer.SvrNameStr + "." + Writer.SuffixStr
}

func (Writer *Writer) CliName() string {
	return Writer.CliPath + "/" + Writer.SvrNameStr + "." + Writer.SuffixStr
}

func (Writer *Writer) ToSvrData() string {
	return ToString(Writer.SvrData)
}

func (Writer *Writer) ToCliData() string {
	return ToString(Writer.SvrData)
}

func ToString(param map[string]any) string {
	byteData, err := json.Marshal(param)
	if err != nil {
		panic(err)
	}
	dataString := string(byteData)
	return dataString
}

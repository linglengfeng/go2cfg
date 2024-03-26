package tojson

import (
	"go2cfg/pkg/tocfgWriter"
)

type Writer struct {
	tocfgWriter.Writer
}

func (Writer *Writer) SvrName() string {
	return Writer.SvrPath + "/" + Writer.SvrNameStr + "." + "json"
}

func (Writer *Writer) CliName() string {
	return Writer.CliPath + "/" + Writer.SvrNameStr + "." + "json"
}

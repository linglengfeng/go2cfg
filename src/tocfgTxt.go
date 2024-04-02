package tocfg

import (
	"os"
)

func writeTxt(writer Writer) error {
	// 修改文件后缀名
	writer.Suffix = ".txt"

	// 创建目录
	createDir(writer.SvrPath)
	createDir(writer.CliPath)

	// server primarykey data
	svrFile := writer.SvrPath + "/" + writer.SvrFileName + writer.Suffix
	svrdataPrimarykey := PrimarykeyJsonStr(writer.PrimarykeySvrData)
	err := os.WriteFile(svrFile, append([]byte(svrdataPrimarykey), byte('\n')), 0644)
	if err != nil {
		return err
	}

	// server unionkey data
	svrUnionFile := writer.SvrPath + "/" + writer.SvrFileName + "_ukey" + writer.Suffix
	svrdataUnionkeys := UnionKeysJsonStr(writer.UnionKeysSvrData)
	err = os.WriteFile(svrUnionFile, append([]byte(svrdataUnionkeys), byte('\n')), 0644)
	if err != nil {
		return err
	}

	// client primarykey data
	cliFile := writer.CliPath + "/" + writer.CliFileName + writer.Suffix
	clidataPrimarykey := PrimarykeyJsonStr(writer.PrimarykeyCliData)
	err = os.WriteFile(cliFile, append([]byte(clidataPrimarykey), byte('\n')), 0644)
	if err != nil {
		return err
	}

	// client unionkey data
	cliUnionFile := writer.CliPath + "/" + writer.CliFileName + "_ukey" + writer.Suffix
	clidataUnionkeys := UnionKeysJsonStr(writer.UnionKeysCliData)
	err = os.WriteFile(cliUnionFile, append([]byte(clidataUnionkeys), byte('\n')), 0644)
	if err != nil {
		return err
	}
	return nil
}

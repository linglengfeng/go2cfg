package tocfg

var (
	globalExportSvr    = prop{Type: EXPORT_SVR, Raw: EXPORT_SVR_RAW, Cell: EXPORT_SVR_CELL, Is: EXPORT_SVR_IS, Name: EXPORT_SVR_NAME}
	globalExportCli    = prop{Type: EXPORT_CLI, Raw: EXPORT_CLI_RAW, Cell: EXPORT_CLI_CELL, Is: EXPORT_CLI_IS, Name: EXPORT_CLI_NAME}
	globalPrimaryKey   = prop{Type: PRIMARY_KEY, Raw: PRIMARY_KEY_RAW, Cell: PRIMARY_KEY_CELL, Name: PRIMARY_KEY_NAME}
	globalUnionSvrKeys = prop{Type: UNION_KEYS_SVR, Raw: UNION_KEYS_SVR_RAW, Cell: UNION_KEYS_SVR_CELL, Name: UNION_KEYS_SVR_NAME}
	globalUnionCliKeys = prop{Type: UNION_KEYS_CLI, Raw: UNION_KEYS_CLI_RAW, Cell: UNION_KEYS_CLI_CELL, Name: UNION_KEYS_CLI_NAME}

	globalOutSvr = prop{Type: OUT_SVR, Raw: OUT_SVR_RAW, Cell: OUT_SVR_CELL}
	globalOutCli = prop{Type: OUT_CLI, Raw: OUT_CLI_RAW, Cell: OUT_CLI_CELL}
	globalType   = prop{Type: TYPE, Raw: TYPE_RAW, Cell: TYPE_CELL}
	globalName   = prop{Type: NAME, Raw: NAME_RAW, Cell: NAME_CELL}
	globalKey    = prop{Type: KEY, Raw: KEY_RAW, Cell: KEY_CELL}
	globalNote   = prop{Type: NOTE, Raw: NOTE_RAW, Cell: NOTE_CELL}

	globalCheckList = []prop{globalExportSvr, globalExportCli, globalPrimaryKey, globalUnionSvrKeys, globalUnionCliKeys, globalOutSvr, globalOutCli, globalType, globalName, globalKey, globalNote}

	globalCheckSameFile  = map[string]struct{}{}
	globalPrimaryKeyInfo = ValInfo{}
)

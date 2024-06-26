package tocfg

// 导出类型
const (
	EXPORT_TYPE_JSON = "json"
	EXPORT_TYPE_TXT  = "txt"
	EXPORT_TYPE_ERL  = "erl"
	EXPORT_TYPE_ERL_HRL  = "erl_hrl"
)

// 格式规范
const (
	EXPORT_SVR      = "EXPORT_SVR"
	EXPORT_SVR_RAW  = 0
	EXPORT_SVR_CELL = 0
	EXPORT_SVR_IS   = 1
	EXPORT_SVR_NAME = 2

	EXPORT_CLI      = "EXPORT_CLI"
	EXPORT_CLI_RAW  = 1
	EXPORT_CLI_CELL = 0
	EXPORT_CLI_IS   = 1
	EXPORT_CLI_NAME = 2

	PRIMARY_KEY      = "PRIMARY_KEY"
	PRIMARY_KEY_RAW  = 2
	PRIMARY_KEY_CELL = 0
	PRIMARY_KEY_NAME = 1

	UNION_KEYS_SVR      = "UNION_KEYS_SVR"
	UNION_KEYS_SVR_RAW  = 3
	UNION_KEYS_SVR_CELL = 0
	UNION_KEYS_SVR_NAME = 1

	UNION_KEYS_CLI      = "UNION_KEYS_CLI"
	UNION_KEYS_CLI_RAW  = 4
	UNION_KEYS_CLI_CELL = 0
	UNION_KEYS_CLI_NAME = 1

	OUT_SVR      = "OUT_SVR"
	OUT_SVR_RAW  = 6
	OUT_SVR_CELL = 0

	OUT_CLI      = "OUT_CLI"
	OUT_CLI_RAW  = 7
	OUT_CLI_CELL = 0

	TYPE      = "TYPE"
	TYPE_RAW  = 8
	TYPE_CELL = 0

	NAME      = "NAME"
	NAME_RAW  = 9
	NAME_CELL = 0

	KEY      = "KEY"
	KEY_RAW  = 10
	KEY_CELL = 0

	NOTE      = "NOTE"
	NOTE_RAW  = 11
	NOTE_CELL = 0

	VALUE      = "VALUE"
	VALUE_CELL = 0
)

// 常量标识
var (
	TRUE  = "TRUE"
	FALSE = "FALSE"
	INT   = "INT"
	STR   = "STR"
	LIST  = "LIST"
)

//默认值
var (
	DEFALUT_EXCELDIR = "./xlsx"
	DEFALUT_SVRDIR   = "./output/svr"
	DEFAULT_CLIDIR   = "./output/cli"
)

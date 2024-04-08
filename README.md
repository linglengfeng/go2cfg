# go2cfg
excel to any file
目前支持 xlsx to txt, xlsx to json, xlsx to erl, xlsx to erl and hrl

## 如何添加别的文件类型
在tocfgWriter.go函数initWorker中添加接口WriterWorker
参照 tocfgTxt.go, tocfgJson.go
格式说明代码：tocfgConst.go

## config.json 说明
		
excel_dir 必填 excel导出目录
writers 必填 导出列表
    type 目前支持 json,txt,erl,erl_hrl
    out_svr_dir 可选，没有字段表示不导出
    out_cli_dir 可选，没有字段表示不导出

导出文件参照：output目录

```json
{
    <!-- 导出目录 -->
    "excel_dir": "./xlsx"
    <!-- 导出列表 -->
    "writers": [
       {
         <!-- 导出格式 json -->
        "type": "json",
        <!-- 导出服务器目录 -->
        "out_svr_dir": "./output/json/svr",
        <!-- 导出客户端目录 -->
        "out_cli_dir": "./output/json/cli"
       },
        <!-- 导出txt -->
       {
        "type": "txt",
        "out_svr_dir": "./output/txt/svr",
        "out_cli_dir": "./output/txt/cli"
       },
       <!-- 导出erl -->
       {
        "type": "erl",
        "out_svr_dir": "./output/erl/svr",
        "out_cli_dir": "./output/erl/cli"
       },
       <!-- 导出erl和头文件 -->
       {
        "type": "erl_hrl",
        "out_svr_dir": "./output/erl_hrl/svr",
        "out_cli_dir": "./output/erl_hrl/cli"
       }
    ]
}
```

### xlsx 格式说明

|   格式说明    |       |       |       |       |       |       |
|---------------|-------|-------|-------|-------|-------|-------|
| EXPORT_SVR	| TRUE	| fName | 			
| EXPORT_CLI	| FALSE	| fName |		
| PRIMARY_KEY	| id	| 		|	
| UNION_KEYS_SVR| 	    |
| UNION_KEYS_CLI|       |       | 
|---------------|-------|-------|-------|-------|-------|-------|
| OUT_SVR	    | TRUE	| TRUE	| TRUE	| TRUE	| TRUE	| FALSE | 
| OUT_CLI	    | TRUE	| TRUE	| TRUE	| TRUE	| TRUE	| TRUE  | 
| TYPE	        | INT	| STR	|  INT	| INT	| LIST	|  STR  | 
| NAME	        | id	| 名字	| 年龄	| 性别	| 物品	| 说明  | 
| KEY	        | id	| name	| age	| sex	| items	| desc  | 
| NOTE	        | 备注	| 备注	| 备注	| 备注	| 备注	|       | 
| VALUE	        | 1	    | name1	| 10	|   1	|[1,2,3]| 说明1 | 
| VALUE	        | 2	    | name2	| 11	|   2	|[1,2,4]| 说明2 | 


第一行：
    第一列：固定 EXPORT_SVR，第二列：TRUE|FALSE 表示该页是否导出服务器（TRUE导出，FALSE不导出），第三列：导出文件名（不填默认filename_sheetname）

第二行：
    第一列：固定 EXPORT_CLI，第二列：TRUE|FALSE 表示该页是否导出客户端（TRUE导出，FALSE不导出），第三列：导出文件名（不填默认filename_sheetname）

第三行：
    第一列：固定 PRIMARY_KEY，第二列：主键名

第四行：
    第一列：固定 UNION_KEYS_SVR，第二列：服务器联合键集合 格式：[[key1],[key1,key2]...]

第五行：
    第一列：固定 UNION_KEYS_CLI，第二列：客户端联合键集合 格式：[[key1],[key1,key2]...]

第六行：
    空行，预留

第七行：
    第一列：固定 OUT_SVR，第x列：TRUE|FALSE 表示该字段是否导出

第八行：
    第一列：固定 OUT_CLI，第x列：TRUE|FALSE 表示该字段是否导出

第九行：
    第一列：固定 TYPE，第x列：该字段类型 INT|STR|LIST

第十行：
    第一列：固定 NAME，第x列：该字段名

第十一行：
    第一列：固定 NOTE，第x列：该字段备注

第x行：
    第一列：固定 VALUE,填VALUE表示该行生效导出，不填或者填其他表示该行不生效不导出

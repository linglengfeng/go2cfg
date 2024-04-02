# go2cfg
excel to file

支持json,txt

想添加别的格式，在tocfgWriter.go 文件中添加case到自定义格式的文件中去执行具体逻辑
参照 tocfgTxt.go, tocfgJson.go

config.json 说明
{
    "excel_dir": "./xlsx"   <!-- 导出目录 -->
    "writers": [ <!-- 导出格式列表 -->
       {
        "type": "json",  <!-- 导出格式 -->
        "out_svr_dir": "./output/json/svr",  <!-- 导出服务器目录 -->
        "out_cli_dir": "./output/json/cli"  <!-- 导出客户端目录 -->
       },
       {
        "type": "txt",
        "out_svr_dir": "./output/txt/svr",
        "out_cli_dir": "./output/txt/cli"
       }
    ]
}
package main

import (
	"flag"
	"fmt"
	"hobby/bin"
	"hobby/ctype"
	"hobby/utils"
)

func showPluginInfo() {
	fmt.Printf(`脚本示例：
	csvbycol2txt >>> {csvbycol2txt xxx.csv 2 xxx.txt} (将fofa.csv中第2列的内容导出（已去重）fofa.txt)
	csvbyname2txt >>> {csvbyname2txt xxx.csv domain xxx.txt} (将fofa.csv中名为domain列的内容导出fofa.txt)
	sleep >>> {sleep 5} (休眠n秒)
	ddcsv >>> {ddcsv ./xxx/ 1k 5} (监测目标文件中是否出现大于1kb的csv文件，超时5s`)
	fmt.Println("\n")

	utils.ShowPlugin()
}

func main() {
	// go build -o hobby.exe -buildmode=exe -ldflags="-s -w" -buildvcs=false -tags=netgo .\main.go
	args := ctype.Args{}
	// Note: 不要在此处解引用标志的返回值
	flushTimeFlag := flag.Int("t", 10, "进程测活 刷新时间")
	hobbyPathFlag := flag.String("c", "go.html", "配置文件地址")
	showPluginInfoFlag := flag.Bool("ph", false, "显示插件信息")
	// 重要: 先解析标志
	flag.Parse()
	if *showPluginInfoFlag {
		showPluginInfo()
		return // Exit program after showing plugin info
	}

	// 现在从已解析的标志中获取值
	args.FlushTime = *flushTimeFlag
	args.HobbyPath = *hobbyPathFlag
	bin.Run(args)

}

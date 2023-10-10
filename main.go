package main

import (
	"flag"
	"fmt"
	"hobby/bin"
	"hobby/ctype"
	"hobby/utils"
)

func showSocketInfo() {
	fmt.Printf("企业端口：\n%v\n", ctype.DScoketPort)
}
func showPluginInfo() {
	fmt.Printf(`脚本示例：
	csvbycol2txt	>>> {csvbycol2txt xxx.csv 2 xxx.txt} (将fofa.csv中第2列的内容导出（已去重）fofa.txt)
	csvbyname2txt	>>> {csvbyname2txt xxx.csv domain xxx.txt} (将fofa.csv中名为domain列的内容导出fofa.txt)
	sleep	>>> {sleep 5} (休眠n秒)
	ddcsv	>>> {ddcsv ./xxx/ 1k 5} (监测目标文件中是否出现大于1kb的csv文件 超时5s)
	logprint	>>> {logprint $ddcsv1$ xxx xxx} (输出脚本变量或字符串到out.txt中以及输出到屏幕上)
	request	>>> {request -u http://xxx/ -f ./xxx.txt [-t] 10 [-head] ./xxx/head [-body] ./xxx/body [-timeout] 10 [-m] GET -o xxx.csv} (http请求 []为可选参数 -m默认GET -t默认10 -u/-f只能选择其一 -timeout默认10)
	socket	>>> {socket -f ./xxx.txt [-p] 1,2,3-10,22,3389 [-t] 10  [-timeout] 10 [-w] ./xxx.txt -o xxx.csv} (socket请求 []为可选参数 -t默认500 -p默认企业端口(-phsocket) -timeout默认10)`)
	fmt.Println("\n")

	utils.ShowPlugin()
}

func main() {
	// go build -o hobby.exe -buildmode=exe -ldflags="-s -w" -buildvcs=false -tags=netgo .\main.go
	args := ctype.Args{}
	// Note: 不要在此处解引用标志的返回值
	flushTimeFlag := flag.Int("t", 10, "进程测活 刷新时间")
	ddprocessFlag := flag.Int("dpt", 120, "监测进程输出 间隔时间")
	hobbyPathFlag := flag.String("c", "go.html", "配置文件地址")
	showPluginInfoFlag := flag.Bool("ph", false, "显示插件信息")
	showPlugSocket := flag.Bool("phsocket", false, "显示插件Socket信息")
	// 重要: 先解析标志
	flag.Parse()
	if *showPluginInfoFlag {
		showPluginInfo()
		return // Exit program after showing plugin info
	}
	if *showPlugSocket {
		showSocketInfo()
		return
	}

	// 现在从已解析的标志中获取值
	args.FlushTime = *flushTimeFlag
	args.HobbyPath = *hobbyPathFlag
	args.Ddprocess = *ddprocessFlag
	bin.Run(args)

}

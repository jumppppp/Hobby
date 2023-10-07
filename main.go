package main

import (
	"flag"
	"fmt"
	"hobby/cplugin"
	"hobby/ctype"
	"hobby/utils"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func run(args ctype.Args) {

	hobby, err := utils.ReadHobby(args.HobbyPath)
	if err != nil {
		fmt.Println("读取失败：", err)
		return
	}
	tags := make([]int, 0, len(hobby))
	for tag := range hobby {
		tags = append(tags, tag)
	}
	sort.Ints(tags)
	wg := &sync.WaitGroup{}
	for _, tag := range tags {
		processes := hobby[tag]
		for _, pn := range processes {

			if pn.Command != "" {
				if pn.Thread > 1 {
					if pn.ThreadContent != "" && pn.ThreadOut != "" {
						Coms, Touts, err := utils.SwapThreadCommand(pn.PPID, pn.Thread, pn.ThreadContent, pn.ThreadOut, pn.Command)
						if err != nil {
							fmt.Printf("[\033[31m进程转换错误\033[0m]{%v} >> %v\n", pn.Command, err)
							return
						}
						var WorkCount int = pn.Thread
						for index, pn := range Coms {
							wg.Add(1)

							go func(wg *sync.WaitGroup, pn string, t int, index int, WorkCount *int) {

								defer func() {
									*WorkCount -= 1
									wg.Done()

								}()
								parts := strings.Fields(pn)

								command := parts[0]
								args := parts[1:]
								pid, err := utils.RunAndGetPID(command, args...)
								if err != nil {
									fmt.Printf("\033[34m(进程%v)\033[0m[\033[31mError\033[0m]{%v} >> %v\n", index, pn, err)
								}
								for {
									info, n := utils.FindProcessByPID(pid)
									if n == 0 {
										fmt.Printf("\033[34m(进程%v)\033[0m[\033[32m执行中...\033[0m]{%v} >> %v\n", index, pn, *info)

									} else if n > -5 {
										fmt.Printf("\033[34m(进程%v)\033[0m[\033[33m执行结束\033[0m]{%v}\n", index, pn)

										break
									} else {
										fmt.Printf("\033[34m(进程%v)\033[0m[\033[31m执行错误\033[0m]{%v} >> %v\n", index, pn, n)
										break
									}
									time.Sleep(time.Second * time.Duration(t))
								}
							}(wg, pn, args.FlushTime, index, &WorkCount)
						}
						wg.Add(1)
						go func(wg *sync.WaitGroup, WorkCount *int, pn ctype.CmdXML, t int) {
							defer wg.Done()
							for {
								if *WorkCount == 0 {
									err = utils.AssembleThreadOut(Touts, pn.ThreadOut)
									if err != nil {
										fmt.Printf("[\033[31m结果聚合失败\033[0m]{%v} >> %v\n", pn.Command, err)
										return
									}
									fmt.Printf("[\033[33m结果聚合完成\033[0m]{%v}\n", pn.Command)
									break

								}
								time.Sleep(time.Second * time.Duration(t))
							}

						}(wg, &WorkCount, pn, args.FlushTime)

					} else {
						fmt.Printf("[\033[31m进程错误\033[0m]{%v} >> 缺失必要值 Thread-Content \n", pn.Command)
						return
					}

				} else {
					wg.Add(1)
					go func(wg *sync.WaitGroup, pn ctype.CmdXML, t int) {

						defer wg.Done()
						parts := strings.Fields(pn.Command)

						command := parts[0]
						args := parts[1:]
						pid, err := utils.RunAndGetPID(command, args...)
						if err != nil {
							fmt.Printf("[\033[31mError\033[0m]{%v} >> %v\n", pn.Command, err)
						}
						for {
							info, n := utils.FindProcessByPID(pid)
							if n == 0 {
								fmt.Printf("[\033[32m执行中...\033[0m]{%v} >> %v\n", pn.Command, *info)

							} else if n > -5 {
								fmt.Printf("[\033[33m执行结束\033[0m]{%v}\n", pn.Command)

								break
							} else {
								fmt.Printf("[\033[31m执行错误\033[0m]{%v} >> %v\n", pn.Command, n)
								break
							}
							time.Sleep(time.Second * time.Duration(t))
						}
					}(wg, pn, args.FlushTime)
				}

			} else if pn.Plugin != "" {
				wg.Add(1)
				go func(wg *sync.WaitGroup, pn ctype.CmdXML, t int) {
					defer wg.Done()

					if pn.Plugin[:1] == "{" && pn.Plugin[len(pn.Plugin)-1:] == "}" {
						pn.Plugin = strings.TrimSpace(pn.Plugin[1 : len(pn.Plugin)-1])
						parts := strings.Fields(pn.Plugin)
						command := parts[0]
						args := parts[1:]
						switch command {

						case "csvbyname2txt":
							_, err := cplugin.ReadCSVbyName(args[0], args[1], args[2])
							if err != nil {
								fmt.Printf("[\033[31m脚本执行错误\033[0m]{%v} >> %v\n", pn.Plugin, err)
								break
							} else {
								fmt.Printf("[\033[33m脚本执行结束\033[0m]{%v}\n", pn.Plugin)
							}
						case "csvbycol2txt":
							column, _ := strconv.Atoi(args[1])
							_, err := cplugin.ReadCSVbyCol(args[0], column, args[2])
							if err != nil {
								fmt.Printf("[\033[31m脚本执行错误\033[0m]{%v} >> %v\n", pn.Plugin, err)
								break
							} else {
								fmt.Printf("[\033[33m脚本执行结束\033[0m]{%v}\n", pn.Plugin)
							}
						default:
							fmt.Printf("[\033[31m脚本不存在\033[0m]{%v}\n", pn.Plugin)
							return
						}
					}
				}(wg, pn, args.FlushTime)

			}

		}
		wg.Wait()
	}
}
func showPluginInfo() {
	fmt.Printf("使用例子：{csvbyname2txt fofa.csv domain fofa.txt}\n解释：将fofa.csv中名为domain列的内容导出为fofa.txt（已去重）\n\n脚本说明：\n原型：csvbycol2txt(fileName string, columnIndex int, outputFileName string)\n原型：csvbyname2txt(fileName string, columnName string, outputFileName string)\n")
}

func main() {
	// go build -o hobby.exe -buildmode=exe -ldflags="-s -w" -buildvcs=false -tags=netgo .\main.go
	args := ctype.Args{}
	// Note: 不要在此处解引用标志的返回值
	flushTimeFlag := flag.Int("t", 2, "进程测活 刷新时间")
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
	run(args)

}

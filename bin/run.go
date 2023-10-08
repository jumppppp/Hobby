package bin

import (
	"hobby/cplugin"
	"hobby/ctype"
	"hobby/utils"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

func Run(args ctype.Args) {
	// 1  读取配置文件
	hobby, err := utils.ReadHobby(args.HobbyPath)
	if err != nil {
		utils.LogPf("[-]读取失败：%v\n", err)
		return
	}
	tags := make([]int, 0, len(hobby))
	for tag := range hobby {
		tags = append(tags, tag)
	}
	sort.Ints(tags)
	// 2 运行链表储存脚本返回值
	InLinkData := make(chan *ctype.LinkData, 1)
	OutLinkData := make(chan *ctype.LinkData, 1)
	Govern := make(chan string, 1)
	go MyGoLink(InLinkData, OutLinkData, Govern)

	// 清理屏幕
	go ClearSrceen(args.FlushTime * 10)
	// 3 依次运行配置文件中的内容
	wg := &sync.WaitGroup{}
	mt := &sync.Mutex{}
	for _, tag := range tags {
		processes := hobby[tag]
		// 循环每个process
		for _, pn := range processes {
			// 外置程序运行
			if pn.Command != "" {
				//多进程
				if pn.Thread > 1 {
					if pn.ThreadContent != "" && pn.ThreadOut != "" {
						Coms, Touts, err := utils.SwapThreadCommand(pn.PPID, pn.Thread, pn.ThreadContent, pn.ThreadOut, pn.Command)
						if err != nil {
							utils.LogPf("[\033[31m进程转换错误\033[0m]{%v} >> %v\n", pn.Command, err)
							return
						}
						var WorkCount int = pn.Thread
						for index, pn := range Coms {
							wg.Add(1)
							// 多进程
							go ManyProcessRun(wg, mt, pn, args.FlushTime, index, &WorkCount, OutLinkData, Govern)
						}
						wg.Add(1)
						// 结果聚合
						go ManyProcessRetCount(wg, &WorkCount, pn, args.FlushTime, Touts)

					} else {
						utils.LogPf("[\033[31m进程错误\033[0m]{%v} >> 缺失必要值 Thread-Content \n", pn.Command)
						return
					}
					// 单进程
				} else {
					wg.Add(1)
					// 单进程
					go OneProcessRun(wg, mt, pn, args.FlushTime, OutLinkData, Govern)
				}
				// 脚本运行
			} else if pn.Plugin != "" {
				wg.Add(1)
				// 脚本进程
				go PluginRun(wg, mt, pn, args.FlushTime, InLinkData, OutLinkData, Govern)

			}

		}
		wg.Wait()
	}
}

// 清理屏幕
func ClearSrceen(num int) {
	for {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
		time.Sleep(time.Duration(num) * time.Second)
	}
}

// 多进程结果聚合
func ManyProcessRetCount(wg *sync.WaitGroup, WorkCount *int, pn ctype.CmdXML, t int, Touts []string) {
	defer wg.Done()
	for {
		if *WorkCount == 0 {
			err := utils.AssembleThreadOut(Touts, pn.ThreadOut)
			if err != nil {
				utils.LogPf("[\033[31m结果聚合失败\033[0m]{%v} >> %v\n", pn.Command, err)
				return
			}
			utils.LogPf("[\033[33m结果聚合完成\033[0m]{%v}\n", pn.Command)
			break

		}
		time.Sleep(time.Second * time.Duration(t))
	}

}

// 多进程执行
func ManyProcessRun(wg *sync.WaitGroup, mt *sync.Mutex, pn string, t int, index int, WorkCount *int, outLinkData chan *ctype.LinkData, govern chan string) {

	defer func() {
		*WorkCount -= 1
		wg.Done()

	}()
	parts := strings.Fields(pn)

	command := parts[0]
	args := parts[1:]

	for i, arg := range args {
		if strings.HasPrefix(arg, "$") && strings.HasSuffix(arg, "$") {
			// 如果字符串以$开始并以$结束，就提取中间的值
			inner := strings.TrimPrefix(arg, "$")
			inner = strings.TrimSuffix(inner, "$")
			mt.Lock()
			govern <- inner
			RetLink := <-outLinkData
			mt.Unlock()
			// 输出或使用替换后的值
			args[i] = RetLink.OkData

		}
	}

	pid, err := utils.RunAndGetPID(command, args...)
	if err != nil {
		utils.LogPf("\033[34m(进程%v)\033[0m[\033[31mError\033[0m]{%v} >> %v\n", index, pn, err)
	}
	for {
		info, n := utils.FindProcessByPID(pid)
		if n == 0 {
			utils.LogPf("\033[34m(进程%v)\033[0m[\033[32m执行中...\033[0m]{%v} >> %v\n", index, pn, *info)

		} else if n > -5 {
			utils.LogPf("\033[34m(进程%v)\033[0m[\033[33m执行结束\033[0m]{%v}\n", index, pn)

			break
		} else {
			utils.LogPf("\033[34m(进程%v)\033[0m[\033[31m执行错误\033[0m]{%v} >> %v\n", index, pn, n)
			break
		}
		time.Sleep(time.Second * time.Duration(t))
	}
}

// 单进程执行
func OneProcessRun(wg *sync.WaitGroup, mt *sync.Mutex, pn ctype.CmdXML, t int, outLinkData chan *ctype.LinkData, govern chan string) {

	defer wg.Done()
	parts := strings.Fields(pn.Command)

	command := parts[0]
	args := parts[1:]
	for i, arg := range args {
		if strings.HasPrefix(arg, "$") && strings.HasSuffix(arg, "$") {
			// 如果字符串以$开始并以$结束，就提取中间的值
			inner := strings.TrimPrefix(arg, "$")
			inner = strings.TrimSuffix(inner, "$")
			mt.Lock()
			govern <- inner
			RetLink := <-outLinkData
			mt.Unlock()
			// 输出或使用替换后的值
			args[i] = RetLink.OkData

		}
	}

	pid, err := utils.RunAndGetPID(command, args...)
	if err != nil {
		utils.LogPf("[\033[31mError\033[0m]{%v} >> %v\n", pn.Command, err)
	}
	for {
		info, n := utils.FindProcessByPID(pid)
		if n == 0 {
			utils.LogPf("[\033[32m执行中...\033[0m]{%v} >> %v\n", pn.Command, *info)

		} else if n > -5 {
			utils.LogPf("[\033[33m执行结束\033[0m]{%v}\n", pn.Command)

			break
		} else {
			utils.LogPf("[\033[31m执行错误\033[0m]{%v} >> %v\n", pn.Command, n)
			break
		}
		time.Sleep(time.Second * time.Duration(t))
	}
}

// 脚本执行
func PluginRun(wg *sync.WaitGroup, mt *sync.Mutex, pn ctype.CmdXML, t int, inLinkData chan *ctype.LinkData, outLinkData chan *ctype.LinkData, govern chan string) {
	defer wg.Done()

	if pn.Plugin[:1] == "{" && pn.Plugin[len(pn.Plugin)-1:] == "}" {
		pn.Plugin = strings.TrimSpace(pn.Plugin[1 : len(pn.Plugin)-1])
		parts := strings.Fields(pn.Plugin)
		command := strings.ToLower(parts[0])
		args := parts[1:]

		for i, arg := range args {
			if strings.HasPrefix(arg, "$") && strings.HasSuffix(arg, "$") {
				// 如果字符串以$开始并以$结束，就提取中间的值
				inner := strings.TrimPrefix(arg, "$")
				inner = strings.TrimSuffix(inner, "$")
				mt.Lock()
				govern <- inner
				RetLink := <-outLinkData
				mt.Unlock()
				// 输出或使用替换后的值
				args[i] = RetLink.OkData

			}
		}

		switch command {
		// 选择脚本
		case "csvbyname2txt":
			_, err := cplugin.ReadCSVbyName(args[0], args[1], args[2])
			if err != nil {
				utils.LogPf("[\033[31m脚本执行错误\033[0m]{%v} >> %v\n", pn.Plugin, err)
			} else {
				utils.LogPf("[\033[33m脚本执行结束\033[0m]{%v}\n", pn.Plugin)
			}
		case "csvbycol2txt":
			column, _ := strconv.Atoi(args[1])
			_, err := cplugin.ReadCSVbyCol(args[0], column, args[2])
			if err != nil {
				utils.LogPf("[\033[31m脚本执行错误\033[0m]{%v} >> %v\n", pn.Plugin, err)
			} else {
				utils.LogPf("[\033[33m脚本执行结束\033[0m]{%v}\n", pn.Plugin)
			}
		case "sleep":
			num, _ := strconv.Atoi(args[0])
			cplugin.CSleep(num)
			utils.LogPf("[\033[33m脚本执行结束\033[0m]{%v}\n", pn.Plugin)
		case "ddcsv":
			num, _ := strconv.Atoi(args[2])
			fname, err := cplugin.MonitorDirCsv(args[0], args[1], num)
			RetLink := ctype.LinkData{UUID: pn.RetMark, OkData: fname}
			inLinkData <- &RetLink
			if err != nil {
				utils.LogPf("[\033[31m脚本执行错误\033[0m]{%v} >> %v\n", pn.Plugin, err)

			} else {
				utils.LogPf("[\033[33m脚本执行结束\033[0m]{%v} >> %v(%v)\n", pn.Plugin, pn.RetMark, fname)
			}
			// test
			govern <- "show"

		case "logprint":
			cplugin.CLogPrint(pn.Plugin, args...)
			utils.LogPf("[\033[33m脚本执行结束\033[0m]{%v}\n", pn.Plugin)
		default:
			utils.LogPf("[\033[31m脚本不存在\033[0m]{%v}\n", pn.Plugin)
			return
		}
	}
}

// 链表保存脚本返回值
func MyGoLink(inLinkData chan *ctype.LinkData, outLinkData chan *ctype.LinkData, govern chan string) {
	LinkT := utils.InitLink()

	for {
		select {
		// 写入数据
		case ldata := <-inLinkData:
			utils.AddRetLink(*ldata, LinkT)
		case control := <-govern:
			switch control {
			case "exit":
				syscall.Exit(1)
			case "show":
				utils.ShowLink(LinkT)
			default:
				tempLink := utils.SelectLinkDatabyUUID(control, LinkT)
				outLinkData <- &tempLink
			}
		}
		time.Sleep(1 * time.Second)
	}
}

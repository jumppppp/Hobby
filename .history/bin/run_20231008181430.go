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

	// 初始化日志文件
	utils.Log_init()
	OutBoardData := ctype.OutBoardData
	KeyBoardDone := ctype.KeyBoardDone
	defer func() {
		close(OutBoardData)
		close(KeyBoardDone)
	}()
	// 键盘 监听
	go cplugin.KeyBoardMain(OutBoardData, KeyBoardDone)
	go cplugin.HandleKeyboardData(OutBoardData)

	// 2 运行链表储存脚本返回值
	InLinkShell := ctype.InLinkShell
	OutLinkShell := ctype.OutLinkShell
	ControlMain := ctype.ControlMain
	InLinkData := ctype.InLinkData
	OutLinkData := ctype.OutLinkData
	Govern := ctype.Govern
	defer func() {
		close(InLinkData)
		close(OutLinkData)
		close(Govern)
		close(InLinkShell)
		close(OutLinkShell)
		close(ControlMain)
	}()
	go LinkShell(InLinkShell, OutLinkShell, ControlMain, InLinkData, OutLinkData, Govern)
	//
	// 清理屏幕
	// go ClearSrceen(args.FlushTime * 10)

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
				PPID := pn.PPID
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
							go ProcessRun(wg, mt, PPID, pn, args.FlushTime, index, &WorkCount, OutLinkData, Govern, args.Ddprocess)
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
					go ProcessRun(wg, mt, PPID, pn, args.FlushTime, -100, nil, OutLinkData, Govern, args.Ddprocess)
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

// 程序执行
func ProcessRun(
	wg *sync.WaitGroup,
	mt *sync.Mutex,
	PPID string,
	pn interface{}, // 可以是string或者ctype.CmdXML
	t int,
	index int, // 对于OneProcessRun为nil
	WorkCount *int, // 对于OneProcessRun为nil
	outLinkData chan *ctype.LinkData,
	govern chan string,
	Dpt int,
) {

	defer wg.Done()
	if WorkCount != nil {
		defer func() {
			*WorkCount -= 1
		}()
	}

	var command string
	var args []string
	var Cmdtext string
	switch v := pn.(type) {
	case string:
		parts := strings.Fields(v)
		command = parts[0]
		args = parts[1:]
		Cmdtext = v
	case ctype.CmdXML:
		parts := strings.Fields(v.Command)
		command = parts[0]
		args = parts[1:]
		Cmdtext = v.Command
	}

	// 处理参数，替换其中的特定值
	for i, arg := range args {
		if strings.HasPrefix(arg, "$") && strings.HasSuffix(arg, "$") {
			inner := strings.TrimPrefix(arg, "$")
			inner = strings.TrimSuffix(inner, "$")
			mt.Lock()
			govern <- inner
			RetLink := <-outLinkData
			mt.Unlock()
			args[i] = RetLink.OkData
		}
	}

	var DDone bool
	DDone = false

	pid, SOut, err := utils.RunAndGetPID(command, args...)
	// 处理程序的输出问题
	go utils.HandelSout(SOut, pid, PPID, &DDone)
	go DdProcessRunStatByLink(PPID, ctype.OutRunStat, Dpt, &DDone)
	go HandleOutRunStat(ctype.OutRunStat, &DDone)
	if err != nil {
		// 打印错误信息
		utils.LogPf("[\033[31mError\033[0m]{%v} >> %v\n", Cmdtext, err)
		return
	}

	for {
		info, n := utils.FindProcessByPID(pid)
		if n == 0 {
			utils.LogPf("\033[34m(进程%v)\033[0m[\033[32m执行中...\033[0m]{%v} >> %v\n", index, Cmdtext, *info)
			// 在运行了，但是输出文件长度一直没有区别

		} else if n > -5 {
			utils.LogPf("\033[34m(进程%v)\033[0m[\033[33m执行结束\033[0m]{%v}\n", index, Cmdtext)
			DDone = true
			break
		} else {
			utils.LogPf("\033[34m(进程%v)\033[0m[\033[31m执行错误\033[0m]{%v} >> %v\n", index, Cmdtext, n)
			DDone = true
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
			// govern <- "show"

		case "logprint":
			cplugin.CLogPrint(pn.Plugin, args...)
			utils.LogPf("[\033[33m脚本执行结束\033[0m]{%v}\n", pn.Plugin)
		default:
			utils.LogPf("[\033[31m脚本不存在\033[0m]{%v}\n", pn.Plugin)
			return
		}
	}
}

// 链表内核
func LinkShell(
	inLink chan *ctype.RetLink,
	outLink chan *ctype.RetLink,
	control chan string,
	inLinkData chan *ctype.LinkData,
	outLinkData chan *ctype.LinkData,
	govern chan string) {

	LinkT := utils.InitLink()
	for {
		select {
		// 写入数据
		case link := <-inLink:
			utils.AddRetLink(link.LinkData, LinkT)

		case c1 := <-control:
			switch control {
			default:
				tempLink := utils.SelectLinkDatabyUUID(c1, LinkT)
				// 监测nil，不能插入空，否则管道阻塞
				if tempLink != nil {
					outLink <- tempLink
				} else {
					outLink <- LinkT
				}

			}
		case ldata := <-inLinkData:
			utils.AddRetLink(*ldata, LinkT)
		case c2 := <-govern:
			switch c2 {
			case "exit":
				syscall.Exit(1)
			case "show":
				utils.ShowLink(LinkT)
			default:
				tempLink := utils.SelectLinkDatabyUUID(c2, LinkT)
				if tempLink != nil {
					outLinkData <- &tempLink.LinkData
				} else {
					outLinkData <- &ctype.LinkData{}
				}

			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// 监测程序运行状态使用链表
func DdProcessRunStatByLink(UUID string, OutRunStat chan *ctype.ProcessRunStat, num int, DDone *bool) {
	ctype.ControlMain <- "HEAD"
	Linkt := <-ctype.OutLinkShell
	for {
		if *DDone {
			return
		}
		tempLink := utils.SelectLinkDatabyUUID(UUID, Linkt)
		// 第一次
		value1, _ := tempLink.LinkData.Data.(int)
		path1 := tempLink.LinkData.OkData
		size1, modTime1, err := utils.GetFileInfo(path1)
		if err != nil {
			// utils.LogPf("[-]DdProcessRunStatByLink Error: %v\n", err)
			continue
		}
		tempRunStat1 := &ctype.ProcessRunStat{
			PID:        value1,
			UUID:       UUID,
			Path:       path1,
			ChangeTime: modTime1.Format("2006-01-02 15:04:05"),
			Length:     int(size1),
		}
		time.Sleep(time.Duration(num) * time.Second)
		// 第二次
		value2, _ := tempLink.LinkData.Data.(int)
		path2 := tempLink.LinkData.OkData
		size2, modTime2, err := utils.GetFileInfo(path1)
		if err != nil {
			// utils.LogPf("[-]DdProcessRunStatByLink Error: %v\n", err)
			continue
		}
		tempRunStat2 := &ctype.ProcessRunStat{
			PID:        value2,
			UUID:       UUID,
			Path:       path2,
			ChangeTime: modTime2.Format("2006-01-02 15:04:05"),
			Length:     int(size2),
		}
		// 发送出去
		OutRunStat <- tempRunStat1
		OutRunStat <- tempRunStat2
		time.Sleep(time.Duration(num) * time.Second)
	}
}
func HandleOutRunStat(OutRunStat chan *ctype.ProcessRunStat, DDone *bool) {
	for {
		if *DDone {
			return
		}
		select {
		case RunStat := <-OutRunStat:
			utils.LogPf("\033[032mPID:[%v]\033[0m\nUUID:[%v]\nPath:[%v]\nModify:[%v]\nLength:[%v]\n\n",
				RunStat.PID, RunStat.UUID, RunStat.Path, RunStat.ChangeTime, RunStat.Length)
		default:

		}
		time.Sleep(time.Second)
	}

}

package cplugin

import (
	"encoding/csv"
	"fmt"
	"hobby/ctype"
	"hobby/utils"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func SocketProbe(target ctype.ProbeTarget, timeout time.Duration, sendData []byte) ctype.ProbeResult {
	address := fmt.Sprintf("%s:%s", target.IP, target.Port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return ctype.ProbeResult{Target: target, IsOpen: false, Err: err, Response: nil}
	}
	defer conn.Close()

	// Send binary data
	if len(sendData) > 0 {
		_, err = conn.Write(sendData)
		if err != nil {
			return ctype.ProbeResult{Target: target, IsOpen: true, Err: err, Response: nil}
		}
	}

	// Read binary response
	response := make([]byte, 1024) // assuming you expect no more than 1024 bytes
	n, err := conn.Read(response)
	if err != nil {
		return ctype.ProbeResult{Target: target, IsOpen: true, Err: err, Response: nil}
	}

	return ctype.ProbeResult{Target: target, IsOpen: true, Err: nil, Response: response[:n]}
}

func GiveSendStream(target ctype.ProbeTarget) map[string][]byte {
	Out := make(map[string][]byte)
	_name := "hobby"
	_domain := "hobby.com"
	Out["SSH"] = ctype.SSH_(_name)
	Out["POP3"] = ctype.FTP_POP3_(_name)
	Out["FTP"] = ctype.FTP_POP3_(_name)
	Out["HTTP"] = ctype.HTTP_(target.IP, target.Port)
	Out["SMTP"] = ctype.SMTP_(_domain)
	Out["NULL"] = ctype.NULL_()
	Out["DNS"] = ctype.DNS_()
	return Out
}

// 返回url列表
func RetIPs(in string) (out []string) {
	lines, err := utils.ReadLines(in)
	if err != nil {
		utils.LogPf(0, "[-]RetIPs err:%v\n", err)
	}

	out = append(out, lines...)

	return out
}

// 返回url列表
func RetStream(in string) []byte {
	lines, err := utils.ReadLines(in)
	if err != nil {
		utils.LogPf(0, "[-]RetStream err:%v\n", err)
	}
	var temp string
	for _, v := range lines {
		temp += v
	}

	return []byte(temp)
}

// 1,2,3,4-10,5,100
func RetPort(in string) (out []string) {
	to := strings.Split(in, ",")
	for _, v := range to {
		too := strings.Split(v, "-")
		if len(too) == 2 {
			min, _ := strconv.Atoi(too[0])
			max, _ := strconv.Atoi(too[1])
			for p := min; p <= max; p++ {
				if p > 0 && p <= 65535 {
					out = append(out, fmt.Sprintf("%v", p))
				}

			}
		} else if len(too) == 1 {
			pp, _ := strconv.Atoi(too[0])
			if pp > 0 && pp <= 65535 {
				out = append(out, fmt.Sprintf("%v", pp))
			}
		} else {
			continue
		}
	}
	return
}

func WriteCSVbySocket(ResqData *ctype.ProbeResult, Outpath string) {
	if ResqData.Err != nil {
		utils.LogPf(2, "Error with the WriteCSVbySocket: %v\n", ResqData.Err)
		return
	}

	file, err := os.OpenFile(Outpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		utils.LogPf(0, "Cannot open file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Extract information from the response and body
	ip := ResqData.Target.IP
	port := ResqData.Target.Port
	response := ResqData.Response
	isopen := ResqData.IsOpen

	// Write data to CSV
	data := []string{"1", ip, port, string(response), strconv.FormatBool(isopen)}
	err = writer.Write(data)
	if err != nil {
		utils.LogPf(0, "Error writing to CSV: %v\n", err)
	}
}

// Handle One
func HandleSocketRespOut(SocketDATA *ctype.SocketToolData) {
	ListResqData := make(map[string]*ctype.ProbeResult)
	for {
		if *SocketDATA.Done {
			if SocketDATA.RetM == "" {
				SocketDATA.RetM = utils.GetUid()
			}
			ReqLink := &ctype.LinkData{UUID: SocketDATA.RetM, OkData: SocketDATA.Out, Data: len(SocketDATA.IPs) * len(SocketDATA.Ports) * len(SocketDATA.SendStream)}
			ctype.InLinkData <- ReqLink
			for _, v := range ListResqData {
				WriteCSVbySocket(v, SocketDATA.Out)
			}
			return
		}
		select {
		case tempResqData := <-SocketDATA.RespOut:
			ListResqData[tempResqData.Target.IP] = tempResqData

		default:

		}
	}
}

// one: http请求响应数据管道 1024
// 处理脚本参数 -> one
func HandleSocketArgs(cmd ctype.CmdXML, SocketDATA *ctype.SocketToolData, args ...string) error {
	defer func() { *SocketDATA.Done = true }()
	SocketDATA.RetM = cmd.RetMark
	SocketDATA.SendStream = make(map[string][]byte)
	for index, value := range args {
		switch value {
		case "-f":
			SocketDATA.IPs = RetIPs(args[index+1])
		case "-p":
			SocketDATA.Ports = RetPort(args[index+1])
		case "-timeout":
			num, _ := strconv.Atoi(args[index+1])
			SocketDATA.TimeOut = time.Duration(num) * time.Second
		case "-o":
			SocketDATA.Out = args[index+1]
		case "-t":
			th, _ := strconv.Atoi(args[index+1])
			SocketDATA.Thread = th
		case "-w":
			SocketDATA.SendStream["me"] = RetStream(args[index+1])
		default:

		}
	}
	if SocketDATA.Out == "" || len(SocketDATA.IPs) == 0 {

		return fmt.Errorf("参数错误：检查[-o|-f]参数")
	}

	// fmt.Printf("%+v\n%+v\n%+v\n", SocketDATA.RespOut, SocketDATA, SocketDATA.Req)
	// 启用批量
	wg := &sync.WaitGroup{}
	// 未指定就变成10
	if SocketDATA.Thread == 0 {
		SocketDATA.Thread = 500
	}
	if len(SocketDATA.Ports) == 0 {
		SocketDATA.Ports = RetPort(ctype.DScoketPort)
	}
	if int(SocketDATA.TimeOut) == 0 {
		SocketDATA.TimeOut = time.Duration(10) * time.Second
	}

	SocketDATA.ThChan = make(chan bool, SocketDATA.Thread)
	//fmt.Printf("%+v\n%+v\n%+v\n", SocketDATA.RespOut, SocketDATA, SocketDATA.Req)
	defer close(SocketDATA.ThChan)
	if len(SocketDATA.IPs) > 0 {
		th_num := 0
		for _, ip := range SocketDATA.IPs {
			SocketDATA.Req.IP = ip
			for _, p := range SocketDATA.Ports {
				SocketDATA.Req.Port = p
				NewStream := GiveSendStream(*SocketDATA.Req)
				for i, v := range NewStream {
					SocketDATA.SendStream[i] = v
				}
				for _, flow := range SocketDATA.SendStream {

					// config := ctype.ProbeTarget{}
					// config.IP = ip
					// config.Port = p
					// url赋值
					// 这里传递过去一个值，在函数中应该又重新映射为一个新的值
					SocketDATA.ThChan <- true
					wg.Add(1)
					if th_num == SocketDATA.Thread-1 {
						th_num = 0
					}
					th_num++
					go ManyRunSocket(*SocketDATA.Req, SocketDATA, wg, flow, th_num, cmd.Plugin)
				}
			}

		}

	}
	wg.Wait()
	return nil

}

// 多个socket执行
func ManyRunSocket(config ctype.ProbeTarget, SocketDATA *ctype.SocketToolData, wg *sync.WaitGroup, flow []byte, index int, cmdtxt string) {
	defer func() {
		<-SocketDATA.ThChan
		wg.Done()
	}()
	// utils.WriteAppendLines([]string{fmt.Sprintf("%v", config)}, "123.txt", true)

	resq := SocketProbe(config, SocketDATA.TimeOut, flow)

	if resq.Err != nil {
		utils.LogPf(0, "\033[032m(线程%v)\033[0m\033[031m[执行错误]\033[0m{%v} >> ERR:%v\n", index, cmdtxt, resq.Err)

	} else {
		SocketDATA.RespOut <- &resq
		utils.LogPf(0, "\033[032m(线程%v)\033[0m\033[033m[执行中]\033[0m{%v} >> IsOpen:%v\n", index, cmdtxt, resq.IsOpen)
	}

}

package cplugin

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"hobby/ctype"
	"hobby/utils"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 用完记得close response
// DoRequest performs a HTTP request based on the given config and returns the response
func DoRequest(client *http.Client, config ctype.RequestConfig, tc chan bool) (*http.Response, []byte, error) {
	// Create a new HTTP request
	req, err := http.NewRequest(config.Method, config.URL, bytes.NewBuffer(config.Body))
	if err != nil {
		return nil, nil, err
	}

	// Set request headers
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	// Set request timeout
	client.Timeout = config.Timeout

	// Perform the request and get the response
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	// Read and return the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}
	return resp, body, nil
}

// 返回头
func RetHeader(in string) map[string]string {
	out := make(map[string]string)
	lines, err := utils.ReadLines(in)
	if err != nil {
		utils.LogPf(0, "[-]RetHeader err:%v\n", err)
	}
	for _, v := range lines {
		out[strings.SplitN(v, ":", 2)[0]] = strings.SplitN(v, ":", 2)[1]
	}

	return out
}

// 返回请求体
func RetBody(in string) []byte {
	lines, err := utils.ReadLines(in)
	if err != nil {
		utils.LogPf(0, "[-]RetBody err:%v\n", err)
	}
	str1 := ""
	for _, v := range lines {
		str1 += v
	}

	return []byte(str1)
}

// 返回url列表
func RetURLs(in string) (out []string) {
	lines, err := utils.ReadLines(in)
	if err != nil {
		utils.LogPf(0, "[-]RetURLs err:%v\n", err)
	}

	out = append(out, lines...)

	return out
}

// 写入csv

func WriteCSV(ResqData *ctype.ResqData, Outpath string) {
	if ResqData.Err != nil {
		utils.LogPf(2, "Error with the request: %v\n", ResqData.Err)
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

	// Write CSV header
	err = writer.Write([]string{"ID", "URL", "IP", "Title", "Status Code"})
	if err != nil {
		utils.LogPf(0, "Error writing to CSV: %v\n", err)
		return
	}

	// Extract information from the response and body
	url := ResqData.PtResq.Request.URL.String()
	ip := ResqData.PtResq.Request.URL.Host // Note: This is the host, not necessarily an IP address.
	statusCode := ResqData.PtResq.StatusCode

	// Extract title from body (simplified, for well-formatted HTML)
	titleStart := strings.Index(string(ResqData.Body), "<title>")
	titleEnd := strings.Index(string(ResqData.Body), "</title>")
	title := "N/A"
	if titleStart != -1 && titleEnd != -1 && titleEnd > titleStart {
		title = string(ResqData.Body[titleStart+7 : titleEnd])
	}

	// Write data to CSV
	data := []string{"1", url, ip, title, fmt.Sprint(statusCode)}
	err = writer.Write(data)
	if err != nil {
		utils.LogPf(0, "Error writing to CSV: %v\n", err)
	}
}

// Handle One
func HandleRespOut(ReqDATA *ctype.RequestToolData) {
	for {
		if *ReqDATA.Done {
			if ReqDATA.RetM == "" {
				ReqDATA.RetM = utils.GetUid()
			}
			ReqLink := &ctype.LinkData{UUID: ReqDATA.RetM, OkData: ReqDATA.Out, Data: len(ReqDATA.Urls)}
			ctype.InLinkData <- ReqLink
			return
		}
		select {
		case tempResqData := <-ReqDATA.RespOut:
			WriteCSV(tempResqData, ReqDATA.Out)
		default:

		}
	}
}

// one: http请求响应数据管道 1024
// 处理脚本参数 -> one
func HandleRequestArgs(cmd ctype.CmdXML, ReqDATA *ctype.RequestToolData, args ...string) error {
	defer func() { *ReqDATA.Done = true }()
	ReqDATA.RetM = cmd.RetMark
	reqArgs := ctype.RequestConfig{Headers: make(map[string]string)}

	ReqDATA.Req = &reqArgs
	for index, value := range args {
		switch value {
		case "-u":
			reqArgs.URL = args[index+1]
		case "-f":
			ReqDATA.Urls = RetURLs(args[index+1])
		case "-head":
			reqArgs.Headers = RetHeader(args[index+1])
		case "-body":
			reqArgs.Body = RetBody(args[index+1])
		case "-m":
			reqArgs.Method = strings.ToUpper(args[index+1])
		case "-timeout":
			num, _ := strconv.Atoi(args[index+1])
			reqArgs.Timeout = time.Duration(num) * time.Second
		case "-o":
			ReqDATA.Out = args[index+1]
		case "-t":
			th, _ := strconv.Atoi(args[index+1])
			ReqDATA.Thread = th
		default:

		}
	}
	if ReqDATA.Out == "" || (ReqDATA.Req.URL == "" && len(ReqDATA.Urls) == 0) {

		return fmt.Errorf("参数错误：检查[-o|-u|-f]参数")
	}

	// fmt.Printf("%+v\n%+v\n%+v\n", ReqDATA.RespOut, ReqDATA, ReqDATA.Req)
	client := &http.Client{}
	// 启用批量
	wg := &sync.WaitGroup{}
	// 未指定就变成10
	if ReqDATA.Thread == 0 {
		ReqDATA.Thread = 10
	}
	if int(ReqDATA.Req.Timeout) == 0 {
		ReqDATA.Req.Timeout = time.Duration(10) * time.Second
	}

	switch ReqDATA.Req.Method {
	case "GET":
	case "HEAD":
	case "POST":
	case "PUT":
	case "PATCH":
	case "DELETE":
	case "CONNECT":
	case "OPTIONS":
	case "TRACE":
	default:
		ReqDATA.Req.Method = "GET"
	}
	ReqDATA.ThChan = make(chan bool, ReqDATA.Thread)
	defer close(ReqDATA.ThChan)
	if len(ReqDATA.Urls) > 0 {
		th_num := 0

		for _, v := range ReqDATA.Urls {
			// url赋值
			ReqDATA.Req.URL = v
			// 这里传递过去一个值，在函数中应该又重新映射为一个新的值
			ReqDATA.ThChan <- true
			wg.Add(1)
			if th_num == ReqDATA.Thread-1 {
				th_num = 0
			}
			th_num++
			go ManyRunReq(client, *ReqDATA.Req, ReqDATA, wg, th_num, cmd.Plugin)

		}

	} else {
		ptresq, body, err := DoRequest(client, *ReqDATA.Req, ReqDATA.ThChan)

		if err != nil {
			utils.LogPf(0, "\033[031m[执行错误]\033[0m{%v} >> ERR:%v\n", cmd.Plugin, err)

		} else {
			TempResq := ctype.ResqData{PtResq: ptresq, Body: body, Err: err}
			ReqDATA.RespOut <- &TempResq
			utils.LogPf(0, "\033[033m[执行中]\033[0m{%v} >> CODE:%v\n", cmd.Plugin, ptresq.StatusCode)
		}
	}
	wg.Wait()
	return nil

}

// 多个url执行
func ManyRunReq(client *http.Client, config ctype.RequestConfig, ReqDATA *ctype.RequestToolData, wg *sync.WaitGroup, index int, cmdtxt string) {
	defer func() {
		<-ReqDATA.ThChan
		wg.Done()
	}()

	ptresq, body, err := DoRequest(client, config, ReqDATA.ThChan)

	if err != nil {
		utils.LogPf(0, "\033[032m(线程%v)\033[0m\033[031m[执行错误]\033[0m{%v} >> ERR:%v\n", index, cmdtxt, err)

	} else {
		TempResq := ctype.ResqData{PtResq: ptresq, Body: body, Err: err}
		ReqDATA.RespOut <- &TempResq
		utils.LogPf(0, "\033[032m(线程%v)\033[0m\033[033m[执行中]\033[0m{%v} >> CODE:%v\n", index, cmdtxt, ptresq.StatusCode)
	}

}

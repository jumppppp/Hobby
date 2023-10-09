package ctype

import (
	"net/http"
	"os"
	"time"
)

type Args struct {
	FlushTime int
	HobbyPath string
	Ddprocess int
}
type ProcessDetails struct {
	PID      int
	Threads  int
	MemoryKB int
}

// ProcessXML 是XML文件的结构表示
type ProcessXML struct {
	Processes []CmdXML `xml:"process"`
}

// CmdXML 是每个进程的结构表示
type CmdXML struct {
	PPID          string
	Tag           int    `xml:"tag"`
	Command       string `xml:"cmd"`
	Plugin        string `xml:"plugin"`
	Thread        int    `xml:"thread"`
	ThreadContent string `xml:"thread-content"`
	ThreadOut     string `xml:"thread-out"`
	RetMark       string `xml:"return-mark"`
}

// 链表的内容
type LinkData struct {
	UUID   string
	Data   interface{}
	OkData string
}

// 链表节点
type RetLink struct {
	LinkData LinkData
	Next     *RetLink
	Prior    *RetLink
}

// 键盘保存
type KeyBoardData struct {
	Rune string
	Key  int
}

// 程序状态
type ProcessRunStat struct {
	PID        int
	PPID       string
	Path       string
	ChangeTime string
	Length     int
}

// linkshell std流
type LinkShellInOutErr struct {
	LinkIn  map[string]*os.File
	LinkOut map[string]*os.File
	LinkErr map[string]*os.File
}

// http单个请求体
type RequestConfig struct {
	URL     string            // URL to send the request to
	Method  string            // HTTP method (GET, POST, etc.)
	Headers map[string]string // HTTP headers
	Body    []byte            // Request body
	Timeout time.Duration     // Request timeout
}

type ResqData struct {
	PtResq *http.Response
	Body   []byte
	Err    error
}

// request工具内容
type RequestToolData struct {
	Req     *RequestConfig
	Urls    []string
	Out     string
	RetM    string
	Thread  int
	ThChan  chan bool
	RespOut chan *ResqData
	Done    *bool
}

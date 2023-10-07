package ctype

type Args struct {
	FlushTime int
	HobbyPath string
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
}

package ctype

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

type LinkData struct {
	UUID   string
	Data   interface{}
	OkData string
}
type RetLink struct {
	LinkData LinkData
	Next     *RetLink
	Prior    *RetLink
}
type KeyBoardData struct {
	Rune string
	Key  int
}
type ProcessRunStat struct {
	PID        int
	PPID       string
	Path       string
	ChangeTime string
	Length     int
}

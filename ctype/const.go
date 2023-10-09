package ctype

var (

	// 键盘记录，缓冲8
	OutBoardData = make(chan *KeyBoardData, 8)
	KeyBoardDone = make(chan bool, 1)

	// 程序状态  一个时间前一个时间后
	OutRunStat = make(chan *ProcessRunStat, 1)

	InLinkData  = make(chan *LinkData, 1)
	OutLinkData = make(chan *LinkData, 1)
	Govern      = make(chan string, 1)

	InLinkShell  = make(chan *RetLink, 1)
	OutLinkShell = make(chan *RetLink, 1)
	ControlMain  = make(chan string, 1)
	// 错误缓存
	ErrLinkShell = make(chan string, 1024)
)

package ctype

var (
	InLinkData  = make(chan *LinkData, 1)
	OutLinkData = make(chan *LinkData, 1)
	Govern      = make(chan string, 1)

	OutBoardData = make(chan *KeyBoardData, 8)
	KeyBoardDone = make(chan bool, 1)
)

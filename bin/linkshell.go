package bin

import (
	"fmt"
	"hobby/ctype"
	"hobby/utils"
	"os"
	"syscall"
	"time"
)

// 链表内核
func LinkShell(
	inLink chan *ctype.RetLink,
	outLink chan *ctype.RetLink,
	control chan string,
	inLinkData chan *ctype.LinkData,
	outLinkData chan *ctype.LinkData,
	govern chan string,
	errLink chan string,

) {

	LinkT := utils.InitLink()
	linkS := &ctype.LinkShellInOutErr{LinkIn: make(map[string]*os.File),
		LinkOut: make(map[string]*os.File),
		LinkErr: make(map[string]*os.File)}
	for {
		select {
		// 写入数据
		case link := <-inLink:
			utils.AddRetLink(link.LinkData, LinkT)

		case c1 := <-control:
			switch control {

			default:
				tempLink := utils.SelectLinkbyUUID(c1, LinkT)
				// 监测nil，不能插入空，否则管道阻塞
				outLink <- tempLink

			}
		case ldata := <-inLinkData:
			utils.AddRetLink(*ldata, LinkT)
		case c2 := <-govern:
			switch c2 {
			case "exit":
				syscall.Exit(999)
			case "show":
				utils.ShowLink(LinkT)
			case "save":
				utils.SaveLink(LinkT)
			case "shell":
				GoLinkShell(linkS)
			case "back":
				BackLinkShell(linkS)
			default:
				tempLink := utils.SelectLinkbyUUID(c2, LinkT)
				outLinkData <- &tempLink.LinkData
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}
func GoLinkShell(linkS *ctype.LinkShellInOutErr) {
	// 保存当前的std
	linkS.LinkOut["out0"] = os.Stdout
	linkS.LinkIn["in0"] = os.Stdin
	linkS.LinkErr["err0"] = os.Stderr
	uid := utils.GetUid()
	utils.WriteCacheByUid(uid, []string{""}, "linkout", false, false)
	utils.WriteCacheByUid(uid, []string{""}, "linkin", false, false)
	utils.WriteCacheByUid(uid, []string{""}, "linkerr", false, false)
	// 创建一个新的文件来作为新的stdout
	newStdout, err := os.OpenFile(fmt.Sprintf("./cache/%v/linkout", uid), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		utils.LogPf("Error creating new stdout:%v\n", err)
		return
	}
	newStdin, err := os.OpenFile(fmt.Sprintf("./cache/%v/linkin", uid), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		utils.LogPf("Error creating new stdout:%v\n", err)
		return
	}
	newStderr, err := os.OpenFile(fmt.Sprintf("./cache/%v/linkerr", uid), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		utils.LogPf("Error creating new stdout:%v\n", err)
		return
	}

	// 将新的stdout设置为标准输出
	os.Stdout = newStdout
	os.Stdin = newStdin
	os.Stderr = newStderr

	linkS.LinkOut["out1"] = newStdout
	linkS.LinkIn["in1"] = newStdin
	linkS.LinkErr["err1"] = newStderr
	utils.LogPf("\033[031m[+]Go LinkShell\033[0m\n")

}
func BackLinkShell(linkS *ctype.LinkShellInOutErr) {
	if _, exists := linkS.LinkOut["out0"]; !exists {
		return
	}
	os.Stdout = linkS.LinkOut["out0"]
	os.Stdin = linkS.LinkOut["in0"]
	os.Stderr = linkS.LinkOut["err0"]
	linkS.LinkOut["out1"].Close()
	linkS.LinkOut["in1"].Close()
	linkS.LinkOut["err1"].Close()
	utils.LogPf("\033[031m[+]Back Shell\033[0m\n")
}

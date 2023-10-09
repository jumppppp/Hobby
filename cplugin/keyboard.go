package cplugin

import (
	"fmt"
	"hobby/ctype"
	"hobby/utils"
	"strings"

	"github.com/eiannone/keyboard"
)

func KeyBoardMain(OutKey chan *ctype.KeyBoardData, done chan bool) {
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	u2 := utils.GetUid()

	defer keyboard.Close()

	// utils.LogPf("[+]Press any key to see the key code. Press Ctrl+Q to quit.\n")

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		// 写入管道

		TempOut := ctype.KeyBoardData{Rune: string(char), Key: int(key)}

		OutKey <- &TempOut
		// 写入文件

		outfile := fmt.Sprintf("[%v|0x%X]", string(char), key)
		err = utils.WriteCacheByUid(u2, []string{outfile}, "KeyBoard.txt", true, true)
		if err != nil {
			utils.LogPf("[-]Error of WriteCacheByUid: %v\n", err)
		}
		// 生成uuid 并且放入linkshell

		u3 := utils.GetUid()

		KeyLinkData := ctype.LinkData{UUID: u3, OkData: outfile}

		ctype.InLinkData <- &KeyLinkData

		if key == keyboard.KeyCtrlQ {
			done <- true
			return
		}

	}

}

// BlinkText 在倒数第二行使用 ANSI 转义码显示闪烁的文本
func BlinkText(text string, delayMillis int) {

	// ANSI 转义码，设置倒数第二行
	// setCursorPosition := fmt.Sprintf("\033[%d;%dH", 1, 0)

	// 清空当前行
	clearLine := "\033[2K"

	// 闪烁的字符串
	blinkingText := "\033[31m" + text + "\033[0m"

	// 将光标移动到倒数第二行
	// fmt.Print(setCursorPosition)

	// 清空当前行
	fmt.Print(clearLine)

	// 打印闪烁的字符串
	fmt.Print(blinkingText)

}

// 处理 监听到的字符
func HandleKeyboardData(OutKey chan *ctype.KeyBoardData) {
	var inputSequence string

	for {
		key := <-OutKey
		inputSequence += key.Rune // Append new input character
		// inputSequence = strings.ToLower(inputSequence)
		// Check length of inputSequence and remove the excess characters at the beginning
		if len(inputSequence) > 100 {
			excess := len(inputSequence) - 100
			inputSequence = inputSequence[excess:]
		}
		BlinkText(inputSequence, 500)
		switch {
		case strings.Contains(inputSequence, "show"):
			ctype.Govern <- "show"
			inputSequence = "" // Reset inputSequence if necessary

		case strings.Contains(inputSequence, "cls"):
			Cls()
			inputSequence = ""

		case strings.Contains(inputSequence, "exit"):
			ctype.Govern <- "exit"
			inputSequence = ""
		case strings.Contains(inputSequence, "shell"):
			ctype.Govern <- "shell"
			inputSequence = ""
		case strings.Contains(inputSequence, "back"):
			ctype.Govern <- "back"
			inputSequence = ""
		case strings.Contains(inputSequence, "save"):
			ctype.Govern <- "save"
			inputSequence = ""
		default:

		}
	}

}

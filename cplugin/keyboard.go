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

		TempOut := ctype.KeyBoardData{Rune: string(char), Key: int(key)}
		OutKey <- &TempOut

		outfile := fmt.Sprintf("[%v|0x%X]", string(char), key)
		err = utils.WriteCacheByUid(u2, []string{outfile}, "KeyBoard.txt", true)
		if err != nil {
			utils.LogPf("[-]Error of WriteCacheByUid: %v\n", err)
		}

		u3 := utils.GetUid()
		KeyLinkData := ctype.LinkData{UUID: u3, OkData: outfile}
		ctype.InLinkData <- &KeyLinkData

		if key == keyboard.KeyCtrlQ {
			done <- true
			return
		}

	}

}
func HandleKeyboardData(OutKey chan *ctype.KeyBoardData) {
	var inputSequence string

	for {
		key := <-OutKey
		inputSequence += key.Rune // Append new input character

		// Check length of inputSequence and remove the excess characters at the beginning
		if len(inputSequence) > 100 {
			excess := len(inputSequence) - 100
			inputSequence = inputSequence[excess:]
		}

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
		default:

		}
	}

}

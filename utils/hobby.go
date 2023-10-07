package utils

import (
	"encoding/xml"
	"fmt"
	"hobby/ctype"
	"os"
	"strings"

	uuid "github.com/satori/go.uuid"
)

func ReadHobby(filename string) (tagMap map[int][]ctype.CmdXML, err error) {

	fileBytes, err := os.ReadFile(filename)
	if err != nil {

		return nil, fmt.Errorf("Failed to read file %s. Error: %v\n", filename, err)
	}
	var p ctype.ProcessXML
	err = xml.Unmarshal(fileBytes, &p)
	if err != nil {

		return nil, fmt.Errorf("Error parsing XML: %v", err)
	}

	tagMap = make(map[int][]ctype.CmdXML)
	for _, process := range p.Processes {
		process.Command = strings.TrimSpace(process.Command)
		process.Plugin = strings.TrimSpace(process.Plugin)
		process.ThreadContent = strings.TrimSpace(process.ThreadContent)
		process.ThreadOut = strings.TrimSpace(process.ThreadOut)
		u1 := uuid.NewV4()
		// 将UUID转换为字符串并获取前16个字符
		process.PPID = u1.String()[:16]
		tagMap[process.Tag] = append(tagMap[process.Tag], process)
	}
	return tagMap, nil
}

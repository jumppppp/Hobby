package utils

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"hobby/ctype"
	"log"
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
		process.PPID = GetUid()
		tagMap[process.Tag] = append(tagMap[process.Tag], process)
	}
	return tagMap, nil
}
func GetUid() string {
	u1 := uuid.NewV4()
	// 将UUID转换为字符串并获取前16个字符
	return u1.String()[:16]
}
func WriteCacheByUid(Uid string, data []string, filename string, _n bool) (err error) {
	folderName := "./cache"

	// create folder
	if _, err = os.Stat(folderName); os.IsNotExist(err) {
		// Folder does not exist, create it
		err = os.Mkdir(folderName, 0755)
		if err != nil {
			// Handle error

			return fmt.Errorf("Failed to create directory: %v\n", err)
		}
	}
	folderName2 := folderName + "/" + Uid
	if _, err = os.Stat(folderName2); os.IsNotExist(err) {
		// Folder does not exist, create it
		err = os.Mkdir(folderName2, 0755)
		if err != nil {
			// Handle error

			return fmt.Errorf("Failed to create directory: %v\n", err)
		}
	}
	path := fmt.Sprintf("%v/%v", folderName2, filename)
	err = WriteAppendLines(data, path, _n)
	if err != nil {
		return fmt.Errorf("Error writing to file: %v - %v", path, err)
	}
	return
}
func WriteAppendLines(lines []string, fileName string, _n bool) error {
	// 使用 os.OpenFile 以追加模式打开文件，如果文件不存在，则创建文件
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		if _n {
			fmt.Fprintln(writer, line)
		} else {
			fmt.Fprintf(writer, "%s", line)
		}
	}
	return writer.Flush()
}
func ShowPlugin() {
	// 指定要分析的包路径
	pkgPath := "cplugin"

	// 使用标准库的go/token包创建一个词法分析器
	fs := token.NewFileSet()

	// 使用标准库的go/parser包解析包的源码
	pkgs, err := parser.ParseDir(fs, pkgPath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	// 遍历包
	for _, pkg := range pkgs {
		// 遍历包的文件
		for _, file := range pkg.Files {
			// 遍历文件中的所有声明
			for i, decl := range file.Decls {
				// 如果是函数声明
				if fn, ok := decl.(*ast.FuncDecl); ok {
					// 打印函数名和类型
					fmt.Printf("\033[032m{%v}[%v]\033[0m \033[033m%s\033[0m \033[031m%s\033[0m\n", file.Name, i, fn.Name.Name, types.ExprString(fn.Type))
				}
			}
		}
	}
}

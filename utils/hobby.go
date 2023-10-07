package utils

import (
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
		u1 := uuid.NewV4()
		// 将UUID转换为字符串并获取前16个字符
		process.PPID = u1.String()[:16]
		tagMap[process.Tag] = append(tagMap[process.Tag], process)
	}
	return tagMap, nil
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
			for _, decl := range file.Decls {
				// 如果是函数声明
				if fn, ok := decl.(*ast.FuncDecl); ok {
					// 打印函数名和类型
					fmt.Printf("Function Name: %s\n", fn.Name.Name)
					fmt.Printf("Function Type: %s\n", types.ExprString(fn.Type))
				}
			}
		}
	}
}

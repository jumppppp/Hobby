package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
)

func main() {
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
					fmt.Printf("{%v}[%v] %s %s\n", file.Name, i, fn.Name.Name, types.ExprString(fn.Type))
				}
			}
		}
	}
}

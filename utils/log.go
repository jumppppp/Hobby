package utils

import (
	"fmt"
	"os"
	"time"
)

func Log_init() {
	// 创建一个日志文件夹（如果不存在）
	logDir := "./result"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	// 创建一个日志文件
	file, err := os.OpenFile("./result/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// fmt.Println("[+]初始化日志")
}

// _flag 0 默认全部记录，1屏幕输出 2日志输出
func LogPf(_flag int, format string, args ...interface{}) {
	text := fmt.Sprintf(format, args...)
	// 获取当前时间
	now := time.Now()
	// 格式化时间
	timen := now.Format("2006-01-02 15:04:05")
	texti := timen + "\t" + text
	if _flag == 1 || _flag == 0 {

		fmt.Printf("%v", text)
	}
	if _flag == 2 || _flag == 0 {
		file, err := os.OpenFile("./result/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		defer func() {
			file.Sync()
			file.Close()

		}()
		_, err = file.WriteString(texti)
		if err != nil {
			panic(err)
		}
	}

}

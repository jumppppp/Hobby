package cplugin

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

func CSleep(num int) {
	time.Sleep(time.Duration(num) * time.Second)
}
func CLogPrint(cmd string, args ...string) {
	// 确保目标目录存在
	dir := "./result"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	// 打开文件，如果不存在则创建
	filePath := filepath.Join(dir, "out.txt")
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 创建多重写入器
	multi := io.MultiWriter(file, os.Stdout)

	// 获取当前时间并格式化
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// 写入日期和参数
	fmt.Fprintf(multi, "{%s}[%v] >> ", currentTime, cmd)
	for _, arg := range args {
		fmt.Fprintf(multi, "%s ", arg)
	}
	fmt.Fprintln(multi) // 添加换行符
}
func MonitorDirCsv(dir string, maxSize string, timeout int) (findName string, err error) {
	var maxSizei int
	var maxSize64 int64
	ftemp := &findName
	if strings.Contains(maxSize, "M") || strings.Contains(maxSize, "m") {
		maxSizei, _ = strconv.Atoi(maxSize[:len(maxSize)-1])
		maxSize64 = int64(maxSizei) * 1024 * 1024
	} else if strings.Contains(maxSize, "K") || strings.Contains(maxSize, "k") {
		maxSizei, _ = strconv.Atoi(maxSize[:len(maxSize)-1])
		maxSize64 = int64(maxSizei) * 1024
	} else {
		return "", fmt.Errorf("脚本运行错误，只能有两种类型，K(k)/M(m)")
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()
	done := make(chan bool, 1)
	go func(done chan bool, ftemp *string, dir string) {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					if strings.HasSuffix(event.Name, ".csv") {
						// 检查文件大小
						fileInfo, err := os.Stat(event.Name)
						if err != nil {
							fmt.Println("Error stating file:", err)
							continue
						}
						if fileInfo.Size() > maxSize64 {
							// fmt.Println("File is larger than limit:", event.Name)
							// 在这里，你可以添加处理大于限制的 CSV 文件的代码
							filename := filepath.Base(event.Name)
							*ftemp = dir + filename
							done <- true
							return
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}(done, ftemp, dir)
	go func(timeout int) {
		time.Sleep(time.Duration(timeout) * time.Second)
		done <- false
	}(timeout)
	// 设置你想要监视的目录
	err = watcher.Add(dir)
	if err != nil {
		return
	}
	ok := <-done
	if ok {
		return findName, nil
	} else {
		return "", fmt.Errorf("监测超时!")
	}

}

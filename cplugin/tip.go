package cplugin

import (
	"fmt"
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

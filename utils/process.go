package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"hobby/ctype"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	uuid "github.com/satori/go.uuid"
)

// 寻找pid值并查询状态
func FindProcessByPID(targetPID int) (*ctype.ProcessDetails, int) {
	cmd := exec.Command("tasklist", "/FO", "CSV", "/NH", "/FI", fmt.Sprintf("PID eq %d", targetPID))
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		//fmt.Errorf("process not found")
		return nil, -1
	}

	lines := strings.Split(out.String(), "\n")
	if len(lines) == 0 {
		//fmt.Errorf("process not found")
		return nil, -1
	}

	// Extracting PID and Memory
	data := strings.Split(lines[0], ",")
	if len(data) < 5 {
		// fmt.Errorf("failed to retrieve process details")
		return nil, -2
	}

	pidStr := strings.Trim(data[1], ` "`)
	memoryStr := strings.Trim(strings.Replace(data[4], " K", "", -1), ` "`)

	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return nil, -9
	}

	memoryKB, err := strconv.Atoi(memoryStr)
	if err != nil {
		return nil, -9
	}

	// Using wmic to get thread count
	out.Reset()
	cmd = exec.Command("wmic", "process", "where", fmt.Sprintf("processid='%d'", pid), "get", "ThreadCount", "/value")
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		return nil, -3
	}

	threadData := strings.Split(out.String(), "=")
	if len(threadData) < 2 {
		//fmt.Errorf("failed to retrieve thread count")
		return nil, -3
	}

	threads, err := strconv.Atoi(strings.TrimSpace(threadData[1]))
	if err != nil {
		return nil, -9
	}

	return &ctype.ProcessDetails{
		PID:      pid,
		Threads:  threads,
		MemoryKB: memoryKB, // converting KB to MB
	}, 0
}

// 执行命令并返回pid
func RunAndGetPID(command string, args ...string) (int, io.ReadCloser, error) {
	cmd := exec.Command(command, args...)
	// 创建一个管道用于读取命令的输出
	output, err := cmd.StdoutPipe()
	if err != nil {
		return 0, nil, err
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // 为了防止命令行窗口的弹出

	err = cmd.Start() // Start command and return immediately without waiting for it to finish
	if err != nil {
		return 0, nil, err
	}

	return cmd.Process.Pid, output, nil
}

// 将多进程的命令进行拆分并返回多条命令
func SwapThreadCommand(PPID string, thread int, tc string, tout string, command string) (Coms []string, Touts []string, err error) {

	// fmt.Println(path)
	lines, err := ReadLines(tc)
	if err != nil {

		return nil, nil, fmt.Errorf("Error reading the file: %v\n", err)
	}
	chunkSize := (len(lines) + thread - 1) / thread
	tcPrefix := filepath.Base(tc)
	tcSuffix := filepath.Ext(tc)
	tcPrefix = tcPrefix[0 : len(tcPrefix)-len(tcSuffix)]

	toPrefix := filepath.Base(tout)
	toSuffix := filepath.Ext(tout)
	toPrefix = toPrefix[0 : len(toPrefix)-len(toSuffix)]
	for i := 0; i < thread; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(lines) {
			end = len(lines)
		}

		wtcFileName := "./cache/" + PPID + "/" + tcPrefix + "_" + strconv.Itoa(i+1) + tcSuffix
		ccname := tcPrefix + "_" + strconv.Itoa(i+1) + tcSuffix
		wtoFileName := "./cache/" + PPID + "/" + toPrefix + "_" + strconv.Itoa(i+1) + toSuffix
		Touts = append(Touts, wtoFileName)
		NewCom := strings.ReplaceAll(command, tc, wtcFileName)
		NewCom = strings.ReplaceAll(NewCom, tout, wtoFileName)
		Coms = append(Coms, NewCom)
		err := WriteCacheByUid(PPID, lines[start:end], ccname, true, false)
		if err != nil {
			return nil, nil, fmt.Errorf("Error writing to file: %v - %v", wtcFileName, err)
		}
	}
	// fmt.Println(Coms, Touts)
	return
}

// 对多进程产生的文件进行聚合
func AssembleThreadOut(touts []string, tout string) (err error) {
	// 创建或打开输出文件
	outputFile, err := os.Create(tout)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)

	// 遍历每个输入文件
	for _, filename := range touts {
		// 打开输入文件
		inputFile, err := os.Open(filename)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(inputFile)
		// 逐行读取并写入到输出文件
		for scanner.Scan() {
			line := scanner.Text()
			writer.WriteString(line + "\n")
		}

		// 检查读取文件过程中是否有错误
		if err := scanner.Err(); err != nil {
			inputFile.Close()
			continue
		}

		inputFile.Close()
	}

	// 确保所有内容都被写入到文件
	writer.Flush()
	return nil
}

// 读取文件
func ReadLines(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// 写入文件
func WriteLines(lines []string, fileName string, _n bool) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func() {
		file.Sync()
		file.Close()

	}()
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
func GetUid() string {
	u1 := uuid.NewV4()
	// 将UUID转换为字符串并获取前16个字符
	return u1.String()[:16]
}

// _n换行，_a追加
func WriteCacheByUid(Uid string, data []string, filename string, _n bool, _a bool) (err error) {
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
	if _a {
		err = WriteAppendLines(data, path, _n)
	} else {
		err = WriteLines(data, path, _n)
	}

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
	defer func() {
		file.Sync()
		file.Close()

	}()
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

// 给出指定文件的大小和修改时间
func GetFileInfo(filename string) (int64, time.Time, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return 0, time.Time{}, err // returning zero value of time.Time if there is an error
	}
	return fileInfo.Size(), fileInfo.ModTime(), nil
}

// 处理程序的输出功能
func HandelSout(Sout io.ReadCloser, newPPID string, DDone *bool) {
	// 读取命令的输出
	buf := make([]byte, 1024)

	// c存入链表
	PPID := strings.Split(newPPID, "#")[0]
	SPid := strings.Split(newPPID, "#")[1]
	outname := fmt.Sprintf("%v.txt", SPid)
	path := strings.ReplaceAll(newPPID, "#"+SPid, "")
	ctype.InLinkData <- &ctype.LinkData{UUID: newPPID, Data: SPid, OkData: fmt.Sprintf("./cache/%v/%v", path, outname)}

	for {
		if *DDone {
			return
		}
		n, err := Sout.Read(buf)
		if err != nil {
			continue
		}
		// 拆解

		if n > 0 {
			err = WriteCacheByUid(PPID, []string{string(buf[:n])}, outname, false, true)
			if err != nil {
				LogPf(0, "Error writing to file: %v - %v\n", outname, err)
			}
		}
	}

}

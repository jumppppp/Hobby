package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"hobby/ctype"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

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

func RunAndGetPID(command string, args ...string) (int, error) {
	cmd := exec.Command(command, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // 为了防止命令行窗口的弹出

	err := cmd.Start() // Start command and return immediately without waiting for it to finish
	if err != nil {
		return 0, err
	}

	return cmd.Process.Pid, nil
}

func SwapThreadCommand(PPID string, thread int, tc string, tout string, command string) (Coms []string, Touts []string, err error) {

	folderName := "./cache"

	// create folder
	if _, err = os.Stat(folderName); os.IsNotExist(err) {
		// Folder does not exist, create it
		err = os.Mkdir(folderName, 0755)
		if err != nil {
			// Handle error

			return nil, nil, fmt.Errorf("Failed to create directory: %v\n", err)
		}
	}
	folderName2 := folderName + "/" + PPID
	if _, err = os.Stat(folderName2); os.IsNotExist(err) {
		// Folder does not exist, create it
		err = os.Mkdir(folderName2, 0755)
		if err != nil {
			// Handle error

			return nil, nil, fmt.Errorf("Failed to create directory: %v\n", err)
		}
	}
	path := fmt.Sprintf("%v/", folderName2)
	// fmt.Println(path)
	lines, err := readLines(tc)
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

		wtcFileName := path + tcPrefix + "_" + strconv.Itoa(i+1) + tcSuffix
		wtoFileName := path + toPrefix + "_" + strconv.Itoa(i+1) + toSuffix
		Touts = append(Touts, wtoFileName)
		NewCom := strings.ReplaceAll(command, tc, wtcFileName)
		NewCom = strings.ReplaceAll(NewCom, tout, wtoFileName)
		Coms = append(Coms, NewCom)
		err := writeLines(lines[start:end], wtcFileName)
		if err != nil {
			return nil, nil, fmt.Errorf("Error writing to file: %v - %v", wtcFileName, err)
		}
	}
	// fmt.Println(Coms, Touts)
	return
}
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
func readLines(fileName string) ([]string, error) {
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

func writeLines(lines []string, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	return writer.Flush()
}

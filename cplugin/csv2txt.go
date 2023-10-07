package cplugin

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
)

func extractAndSaveValues(reader *csv.Reader, colIndex int, outputFileName string) ([]string, error) {
	uniqueValues := make(map[string]struct{})

	for {
		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, csv.ErrFieldCount) {
				continue
			}
			break
		}

		value := record[colIndex]
		uniqueValues[value] = struct{}{}
	}

	uniqueList := make([]string, 0, len(uniqueValues))
	for value := range uniqueValues {
		uniqueList = append(uniqueList, value)
	}

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return nil, err
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	for _, value := range uniqueList {
		if err := writer.Write([]string{value}); err != nil {
			return nil, err
		}
	}
	writer.Flush()

	return uniqueList, nil
}

func ReadCSVbyCol(fileName string, columnIndex int, outputFileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}
	if columnIndex < 0 || columnIndex >= len(headers) {
		return nil, fmt.Errorf("列索引 %d 超出文件列数范围", columnIndex)
	}

	return extractAndSaveValues(reader, columnIndex, outputFileName)
}

func ReadCSVbyName(fileName string, columnName string, outputFileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	var colIndex int = -1
	for i, header := range headers {
		if header == columnName {
			colIndex = i
			break
		}
	}

	if colIndex == -1 {
		return nil, fmt.Errorf("列名 %s 不存在于文件中", columnName)
	}

	return extractAndSaveValues(reader, colIndex, outputFileName)
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

type filehandle struct {
}

var FileHandle = &filehandle{}

// ReadLines 读取文件
func (f *filehandle) ReadLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// WriteLines 写入文件
func (f *filehandle) WriteLines(filename string, lines []string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		if line == "" {
			continue
		}
		fmt.Fprintln(writer, line)
	}
	return writer.Flush()
}

// UniqueStrings 去重
func (f *filehandle) UniqueStrings(slice []string) []string {
	seen := make(map[string]struct{}) // 使用空结构体来节省内存
	result := []string{}
	for _, value := range slice {
		if _, exists := seen[value]; !exists {
			seen[value] = struct{}{}
			result = append(result, value)
		}
	}
	return result
}

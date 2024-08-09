package utils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

type zipHandle struct{}

var ZipHandle = &zipHandle{}

// CreateZipStream 创建zip文件流
func (z *zipHandle) CreateZipStream(name, content string) ([]byte, error) {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	fileHeader := &zip.FileHeader{Name: name, Method: zip.Deflate}
	fileWriter, err := zipWriter.CreateHeader(fileHeader)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(fileWriter, bytes.NewReader([]byte(content)))
	if err != nil {
		return nil, err
	}

	err = zipWriter.Flush()
	if err != nil {
		return nil, err
	}

	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// CreateZipInMemory 创建zip文件流
func (z *zipHandle) CreateZipInMemory(fileName string, Content string) ([]byte, error) {
	var buffer bytes.Buffer
	zipWriter := zip.NewWriter(&buffer)
	fileContent := []byte(Content)

	fileHeader := &zip.FileHeader{
		Name:          fileName,    // 文件名
		Method:        zip.Deflate, // 压缩方法
		Modified:      time.Now(),
		ExternalAttrs: 0644,
	}
	writer, err := zipWriter.CreateHeader(fileHeader)
	if err != nil {
		return nil, err
	}
	_, err = writer.Write(fileContent)
	if err != nil {
		return nil, err
	}
	err = zipWriter.Close()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// readZipFileFromLocal 从本地读取zip的文件并返回流
func (z *zipHandle) readZipFileFromLocal(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening ZIP file:", err)
		return
	}
	defer file.Close()
	// 获取文件信息
	info, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}
	// 创建ZIP读取器
	reader, err := zip.NewReader(file, info.Size())
	if err != nil {
		fmt.Println("Error creating ZIP reader:", err)
		return
	}

	// 遍历ZIP文件中的所有文件头
	for _, fileHeader := range reader.File {
		fmt.Println("Name:", fileHeader.Name)
		fmt.Println("Method:", fileHeader.Method)
		fmt.Println("Modified:", fileHeader.Modified.Format(time.RFC1123))
		fmt.Println("CRC32:", fileHeader.CRC32)
		fmt.Println("Compressed Size:", fileHeader.CompressedSize)
		fmt.Println("Uncompressed Size:", fileHeader.UncompressedSize)
		fmt.Println("Extra:", string(fileHeader.Extra))
		fmt.Println("ExternalAttrs:", fileHeader.ExternalAttrs)
		fmt.Println("Comment:", fileHeader.Comment)
		fmt.Println("------------------------")
	}
}

// 从zip流读取到zip里的文件名和文件内容 readFileFromZipStream
func (z *zipHandle) readFileFromZipStream(content []byte) (string, error) {
	r := bytes.NewReader(content)
	// 打开ZIP文件
	zr, err := zip.NewReader(r, int64(len(content)))
	if err != nil {
		return "", nil
	}
	allContent := ""
	for _, f := range zr.File {
		fmt.Println("file name:", f.Name)
		// 打开文件
		rc, err := f.Open()
		if err != nil {
			return "", nil
		}
		defer rc.Close()
		// 读取文件内容
		data, err := io.ReadAll(rc)
		if err != nil {
			return "", nil
		}
		allContent += string(data)
	}
	return allContent, nil
}

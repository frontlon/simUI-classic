package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

//读取csv
func ReadCsv(filename string) ([][]string, error) {

	if filename == "" {
		return nil, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	data := [][]string{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}
		data = append(data, record)
	}
	return data, nil
}

//写入csv
func WriteCsv(filename string, data [][]string) error {

	if filename == "" {
		return nil
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("open file is failed, err: ", err)
		return err
	}
	defer file.Close()
	file.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(file)
	for _, d := range data {
		fmt.Println("write",len(d),d)
		w.Write(d)
	}
	w.Flush()
	return nil
}

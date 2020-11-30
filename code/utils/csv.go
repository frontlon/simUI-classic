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
		}
		data = append(data, record)
	}
	return data, nil
}
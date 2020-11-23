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
		create := []string{"","","","","","",""}
		if err == io.EOF {
			break
		}else{
			create = record;
		}
		data = append(data, create)
	}
	return data, nil
}
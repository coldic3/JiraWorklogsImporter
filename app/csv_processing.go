package app

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func ReadCSVFile(filePath string) ([][]string, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(fmt.Sprintf("%s/%s", dir, filePath))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	return reader.ReadAll()
}

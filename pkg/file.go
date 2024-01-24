package pkg

import (
	"encoding/csv"
	"io"
	"os"
)

func ReadDataFromFile(inputFile string, creator Creator) ([]interface{}, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ReadDataFromReader(file, creator)
}

func ReadDataFromReader(reader io.Reader, creator Creator) ([]interface{}, error) {
	csvReader := csv.NewReader(reader)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	inputs := make([]interface{}, len(data))
	for i, v := range data {
		inputs[i] = creator.ConvertToInput(v)
	}

	return inputs, nil
}

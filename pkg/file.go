package pkg

import (
	"encoding/csv"
	"os"
)

func ReadDataFromFile(inputFile string, creator Creator) ([]interface{}, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	inputs := make([]interface{}, len(data))
	for i, v := range data {
		inputs[i] = creator.ConvertToInput(v)
	}

	return inputs, nil
}

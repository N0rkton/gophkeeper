// Package verysecretfiles provides functions for reading and writing data to a JSON file.
package verysecretfiles

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"

	"gophkeeper/internal/datamodels"
)

// ReadData reads data from a JSON file and returns a map of datamodels.UniqueData to datamodels.Data.
func ReadData() (map[datamodels.UniqueData]datamodels.Data, error) {
	file, err := os.OpenFile("data.json", os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, errors.New("failed to open file")
	}
	defer file.Close()

	store := make(map[datamodels.UniqueData]datamodels.Data)
	var data []datamodels.Data
	var tmp datamodels.Data

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		err = json.Unmarshal(scanner.Bytes(), &tmp)
		if err != nil {
			if err.Error() != "EOF" {
				return nil, errors.New("failed to decode data")
			}
		}
		data = append(data, tmp)
	}

	for _, v := range data {
		store[datamodels.UniqueData{DataID: v.DataID, UserID: v.UserID}] = datamodels.Data{
			DataID:    v.DataID,
			UserID:    v.UserID,
			Data:      v.Data,
			Deleted:   v.Deleted,
			Metadata:  v.Metadata,
			ChangedAt: v.ChangedAt,
		}
	}

	return store, nil
}

// WriteData writes the provided data to a JSON file.
func WriteData(data datamodels.Data) error {
	file, err := os.OpenFile("data.json", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return errors.New("failed to open file")
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		return errors.New("failed to encode data")
	}

	return nil
}

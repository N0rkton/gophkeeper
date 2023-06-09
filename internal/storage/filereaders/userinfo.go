// Package filereaders provides functions for reading and writing data to a JSON file.
package filereaders

import (
	"bufio"
	"encoding/json"
	"errors"
	"gophkeeper/internal/datamodels"
	"gophkeeper/internal/sessionstorage"

	"os"
)

// ReadUsers reads data from a JSON file and returns sessionstorage.UserSession.
func ReadUsers() (sessionstorage.UserSession, error) {
	file, err := os.OpenFile("users.json", os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return sessionstorage.UserSession{}, errors.New("failed to open file")
	}
	defer file.Close()
	user := sessionstorage.Init()
	var data []datamodels.Auth
	var tmp datamodels.Auth
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		err = json.Unmarshal(scanner.Bytes(), &tmp)
		if err != nil {
			if err.Error() != "EOF" {
				return sessionstorage.UserSession{}, errors.New("failed to decode data")
			}
		}
		data = append(data, tmp)
	}
	for _, v := range data {
		user.AddUser(v.Login, v.Password, v.ID)
	}
	return user, nil
}

// WriteUser writes the provided data to a JSON file.
func WriteUser(auth datamodels.Auth) error {
	file, err := os.OpenFile("users.json", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return errors.New("failed to open file")
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(auth)
	if err != nil {
		return errors.New("failed to encode data")
	}
	return nil
}

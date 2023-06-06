package verysecretfiles

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"gophkeeper/internal/sessionstorage"
	"log"
	"os"
)

func Init() (sessionstorage.SessionStorage, error) {
	file, err := os.OpenFile("users.txt", os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatalf("users file not found")
	}
	defer file.Close()
	var text []string
	usersStorage := sessionstorage.NewAuthUsersStorage()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		err = json.Unmarshal(scanner.Bytes(), &text)
		if err != nil {
			return usersStorage,
				fmt.Errorf("unable to unmarshall: %w", err)
		}
		//id, err := strconv.Atoi(text[2])
		if err != nil {
			return nil, errors.New("wrong data in file")
		}
		//usersStorage.AddUserFromFile(text[0], text[1], uint32(id))

	}
	return usersStorage, nil
}

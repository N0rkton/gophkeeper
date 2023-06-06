package sessionstorage

import (
	"errors"
	"sync"
)

// todo testi i doc
type SessionStorage interface {
	AddUser(user string, id uint32) error

	GetUser(user string) (uint32, error)
}
type authUsersStorage struct {
	authUsers map[string]uint32
	mutex     sync.RWMutex
}

func NewAuthUsersStorage() SessionStorage {
	return &authUsersStorage{authUsers: make(map[string]uint32)}
}

func (us *authUsersStorage) AddUser(user string, id uint32) error {
	us.mutex.Lock()
	us.authUsers[user] = id
	us.mutex.Unlock()
	return nil
}

func (us *authUsersStorage) GetUser(user string) (uint32, error) {
	us.mutex.RLock()
	id, ok := us.authUsers[user]
	us.mutex.RUnlock()
	if !ok {
		return 0, errors.New("user not found")
	}
	return id, nil
}

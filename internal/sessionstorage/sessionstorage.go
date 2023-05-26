package sessionstorage

import (
	"errors"
	"sync"
	"sync/atomic"
)

type SessionStorage interface {
	AddUser(login string, password string) error
	AddUserFromFile(login string, password string, id uint32) error
	GetUser(login string) (User, error)
}
type authUsersStorage struct {
	authUsers map[string]User
	mutex     sync.RWMutex
	lastID    uint32
}
type User struct {
	Password string
	Id       uint32
}

func NewAuthUsersStorage() SessionStorage {
	return &authUsersStorage{authUsers: make(map[string]User)}
}
func (us *authUsersStorage) AddUser(login string, password string) error {
	atomic.AddUint32(&us.lastID, 1)
	us.mutex.Lock()
	us.authUsers[login] = User{Password: password, Id: us.lastID}
	us.mutex.Unlock()
	return nil
}
func (us *authUsersStorage) AddUserFromFile(login string, password string, id uint32) error {
	us.lastID = id
	us.mutex.Lock()
	us.authUsers[login] = User{Password: password, Id: us.lastID}
	us.mutex.Unlock()
	return nil
}
func (us *authUsersStorage) GetUser(login string) (User, error) {
	us.mutex.RLock()
	user, ok := us.authUsers[login]
	us.mutex.RUnlock()

	if !ok {
		return User{}, errors.New("user not found")
	}
	return user, nil
}

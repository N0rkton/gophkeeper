// Package sessionstorage provides an implementation of SessionStorage for storing user session information.
package sessionstorage

import (
	"errors"
	"sync"
)

// SessionStorage defines the methods for managing user sessions.
type SessionStorage interface {
	AddUser(user string, id uint32) error

	GetUser(user string) (uint32, error)
}

// authUsersStorage is an implementation of SessionStorage that stores user session data in memory.
type authUsersStorage struct {
	authUsers map[string]uint32
	mutex     sync.RWMutex
}

// NewAuthUsersStorage creates a new instance of authUsersStorage.
func NewAuthUsersStorage() SessionStorage {
	return &authUsersStorage{authUsers: make(map[string]uint32)}
}

// AddUser adds a new user to the session storage.
func (us *authUsersStorage) AddUser(user string, id uint32) error {
	us.mutex.Lock()
	us.authUsers[user] = id
	us.mutex.Unlock()
	return nil
}

// GetUser retrieves the user ID from the session storage based on the username.
func (us *authUsersStorage) GetUser(user string) (uint32, error) {
	us.mutex.RLock()
	id, ok := us.authUsers[user]
	us.mutex.RUnlock()
	if !ok {
		return 0, errors.New("user not found")
	}
	return id, nil
}

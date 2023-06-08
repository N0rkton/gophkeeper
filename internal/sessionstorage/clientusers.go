// Package sessionstorage provides an in-memory session storage implementation.
package sessionstorage

import (
	"errors"
	"gophkeeper/internal/datamodels"
)

// UserSession represents a session storage that holds user information.
type UserSession struct {
	users map[string]datamodels.Login
}

// Init initializes a new UserSession instance.
func Init() UserSession {
	return UserSession{make(map[string]datamodels.Login)}
}

// AddUser adds a new user to the session storage.
func (u *UserSession) AddUser(login string, password string, id uint32) error {
	_, ok := u.GetUser(login)
	if ok {
		return errors.New("user already exists")
	}
	u.users[login] = datamodels.Login{Password: password, ID: id}
	return nil
}

// GetUser retrieves a user from the session storage based on the login.
// It returns the user and a boolean indicating if the user exists.
func (u *UserSession) GetUser(login string) (datamodels.Login, bool) {
	user, ok := u.users[login]
	if !ok {
		return datamodels.Login{}, ok
	}
	return user, ok
}

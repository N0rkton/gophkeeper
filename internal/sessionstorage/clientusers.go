package sessionstorage

import "errors"

type UserSession struct {
	users map[string]string
}

func Init() UserSession {
	return UserSession{make(map[string]string)}
}
func (u *UserSession) AddUser(user string, password string) error {
	_, ok := u.GetUser(user)
	if ok {
		return errors.New("user already exists")
	}
	u.users[user] = password
	return nil
}

func (u *UserSession) GetUser(user string) (string, bool) {
	password, ok := u.users[user]
	if !ok {
		return "", ok
	}
	return password, ok
}

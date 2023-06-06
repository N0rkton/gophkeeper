package sessionstorage

import "errors"

type User struct {
	Password string
	Id       uint32
}
type UserSession struct {
	users map[string]User
}

func Init() UserSession {
	return UserSession{make(map[string]User)}
}
func (u *UserSession) AddUser(login string, password string, id uint32) error {
	_, ok := u.GetUser(login)
	if ok {
		return errors.New("user already exists")
	}
	u.users[login] = User{Password: password, Id: id}
	return nil
}

func (u *UserSession) GetUser(login string) (User, bool) {
	user, ok := u.users[login]
	if !ok {
		return user, ok
	}
	return user, ok
}

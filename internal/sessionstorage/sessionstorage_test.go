package sessionstorage

import (
	"fmt"
	"log"
)

func ExampleUserSession_AddUser() {
	user := Init()
	err := user.AddUser("login", "password", 0)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(user.users["login"].Password)
	fmt.Println(user.users["login"].ID)
	//Output:
	//password
	//0
}
func ExampleUserSession_GetUser() {
	user := Init()
	err := user.AddUser("login", "password", 0)
	if err != nil {
		log.Fatalln(err)
	}
	u, ok := user.GetUser("login")
	if !ok {
		log.Fatalln(ok)
	}
	fmt.Println(u.Password)
	fmt.Println(u.ID)
	//Output:
	//password
	//0
}
func ExampleAuthUsersStorage_AddUser_GetUser() {
	user := NewAuthUsersStorage()
	err := user.AddUser("userToken", 0)
	if err != nil {
		log.Fatalln(err)
	}
	id, err := user.GetUser("userToken")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(id)
	//Output:
	//0
}

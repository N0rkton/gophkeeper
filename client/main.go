package main

import (
	"gophkeeper/internal/storage"
)

func Init() {

	Storage := storage.NewMemoryStorage()

}
func main() {
	Init()
	//TODO как делать запросы на клиенте
	//TODO воркер
	//TODO
}

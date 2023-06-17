package main

import (
	"log"
	"os"

	"gophkeeper/internal/actions"
	"gophkeeper/internal/storage"

	"github.com/urfave/cli/v2"
)

func Init() storage.Storage {
	storage.Init()
	return storage.NewMemoryStorage()
}
func main() {
	store := Init()

	app := cli.NewApp()
	app.Name = "password keeper"
	app.Usage = "keeps your passwords"
	app.Description = "GophKeeper представляет собой клиент-серверную систему, позволяющую пользователю надёжно и безопасно хранить логины, пароли, бинарные данные и прочую приватную информацию."
	app.Action = actions.MainAction

	app.Commands = []*cli.Command{

		actions.Auth(store),
		actions.GetData(store),
		actions.AddData(store),
		actions.Sync(store),
		actions.DelData(store),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}

}

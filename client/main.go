package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"gophkeeper/internal/datamodels"
	"gophkeeper/internal/storage"
	"log"
	"os"
)

var Storage storage.Storage

func Init() {
	Storage = storage.NewMemoryStorage()
	storage.Init()
}
func main() {
	Init()
	app := cli.NewApp()
	app.Name = "password keeper"
	app.Usage = "keeps your passwords"
	app.Description = "GophKeeper представляет собой клиент-серверную систему, позволяющую пользователю надёжно и безопасно хранить логины, пароли, бинарные данные и прочую приватную информацию."
	app.Action = mainAction

	app.Commands = []*cli.Command{

		auth(),
		getData(),
		addData(),
		sync(),
		delData(),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}

}
func auth() *cli.Command {
	return &cli.Command{
		Name:    "authentication",
		Usage:   "used to authenticate new users; you need to enter login and password",
		Aliases: []string{"a", "auth"},
		Action: func(ctx *cli.Context) error {
			n := ctx.NArg()
			if n == 0 {
				return fmt.Errorf("no argument provided for auth")
			}
			if n != 2 {
				return fmt.Errorf("not enough arguments provided for auth")
			}
			login := ctx.Args().Get(0)
			password := ctx.Args().Get(1)
			err := Storage.Auth(login, password)
			if err != nil {
				return fmt.Errorf("error happend: %w", err)
			}
			fmt.Println("successful auth")
			return nil
		},
	}
}
func addData() *cli.Command {
	return &cli.Command{
		Name:    "addData",
		Usage:   "used to add new data to keep it; you need to enter login and password, then data name, data and meta information if needed",
		Aliases: []string{"add"},
		Action: func(ctx *cli.Context) error {
			n := ctx.NArg()
			if n == 0 {
				return fmt.Errorf("no argument provided for auth")
			}
			if n < 4 {
				return fmt.Errorf("not enough arguments provided for auth")
			}
			login := ctx.Args().Get(0)
			password := ctx.Args().Get(1)
			id, err := Storage.Login(login, password)
			if err != nil {
				return fmt.Errorf("error login happend: %w", err)
			}
			var data datamodels.Data
			data.DataID = ctx.Args().Get(2)
			data.Data = ctx.Args().Get(3)
			data.Metadata = ctx.Args().Get(4)
			data.UserID = id
			err = Storage.AddData(data)
			if err != nil {
				return fmt.Errorf("error add happend: %w", err)
			}
			fmt.Println("data added successfully")
			return nil
		},
	}
}
func getData() *cli.Command {
	return &cli.Command{
		Name:    "get data",
		Usage:   "used to get data ; you need to enter login and password, then data name",
		Aliases: []string{"get", "g"},
		Action: func(ctx *cli.Context) error {
			n := ctx.NArg()
			if n == 0 {
				return fmt.Errorf("no argument provided for auth")
			}
			if n != 3 {
				return fmt.Errorf("wrong amount of arguments")
			}
			login := ctx.Args().Get(0)
			password := ctx.Args().Get(1)
			id, err := Storage.Login(login, password)
			if err != nil {
				return fmt.Errorf("error login happend: %w", err)
			}
			dataId := ctx.Args().Get(2)
			data, err := Storage.GetData(dataId, id)
			if err != nil {
				return fmt.Errorf("error get happend: %w", err)
			}
			fmt.Println("DataID: " + data.DataID + " Data: " + data.Data + " Meta Info: " + data.Metadata)
			return nil
		},
	}
}
func delData() *cli.Command {
	return &cli.Command{
		Name:    "Delete data",
		Usage:   "used to delete data ; you need to enter login and password, then data name",
		Aliases: []string{"del", "d"},
		Action: func(ctx *cli.Context) error {
			n := ctx.NArg()
			if n == 0 {
				return fmt.Errorf("no argument provided for auth")
			}
			if n != 3 {
				return fmt.Errorf("wrong amount of arguments")
			}
			login := ctx.Args().Get(0)
			password := ctx.Args().Get(1)
			id, err := Storage.Login(login, password)
			if err != nil {
				return fmt.Errorf("error login happend: %w", err)
			}
			dataId := ctx.Args().Get(2)
			err = Storage.DelData(dataId, id)
			if err != nil {
				return fmt.Errorf("error deletr happend: %w", err)
			}
			fmt.Println("data deleted successfully")
			return nil
		},
	}
}
func sync() *cli.Command {
	return &cli.Command{
		Name:    "synchronization",
		Usage:   "used synchronize server and client; you need to enter login and password",
		Aliases: []string{"sync", "s"},
		Action: func(ctx *cli.Context) error {
			n := ctx.NArg()
			if n == 0 {
				return fmt.Errorf("no argument provided for auth")
			}
			if n != 2 {
				return fmt.Errorf("wrong amount of arguments")
			}
			login := ctx.Args().Get(0)
			password := ctx.Args().Get(1)
			id, err := Storage.Login(login, password)
			if err != nil {
				return fmt.Errorf("error login happend: %w", err)
			}
			err = Storage.ClientSync(id, nil)
			if err != nil {
				return fmt.Errorf("error client sync happend: %w", err)
			}
			data, err := Storage.Sync(id)
			if err != nil {
				return fmt.Errorf("error sync happend: %w", err)
			}
			for _, v := range data {
				fmt.Println("DataID: " + v.DataID + " Data: " + v.Data + " Meta Info: " + v.Metadata)
			}

			return nil
		},
	}
}
func mainAction(ctx *cli.Context) error {
	ctx.App.Command("help").Run(ctx)
	return nil
}

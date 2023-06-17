package actions

import (
	"fmt"
	"gophkeeper/internal/datamodels"
	"gophkeeper/internal/storage"

	"github.com/urfave/cli/v2"
)

// Auth - used to authenticate new users;
func Auth(store storage.Storage) *cli.Command {
	return &cli.Command{
		Name:  "authentication",
		Usage: "used to authenticate new users; you need to enter login and password; example: go run main.go auth login password",

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
			err := store.Auth(login, password)
			if err != nil {
				return fmt.Errorf("error happend: %w", err)
			}
			fmt.Println("successful auth")
			return nil
		},
	}
}

// AddData - used to add new data to keep it
func AddData(store storage.Storage) *cli.Command {
	return &cli.Command{
		Name:    "addData",
		Usage:   "used to add new data to keep it; you need to enter login and password, then data name, data and meta information if needed; example: go run main.go add login password dataID data metaData",
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

			id, err := store.Login(login, password)
			if err != nil {
				return fmt.Errorf("error login happend: %w", err)
			}
			var data datamodels.Data
			data.DataID = ctx.Args().Get(2)
			data.Data = ctx.Args().Get(3)
			data.Metadata = ctx.Args().Get(4)
			data.UserID = id
			err = store.AddData(data)
			if err != nil {
				return fmt.Errorf("error add happend: %w", err)
			}
			fmt.Println("data added successfully")
			return nil
		},
	}
}

// GetData - used to get data
func GetData(store storage.Storage) *cli.Command {
	return &cli.Command{
		Name:    "get data",
		Usage:   "used to get data ; you need to enter login and password, then data name; example: go run main.go get login password dataId",
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
			id, err := store.Login(login, password)
			if err != nil {
				return fmt.Errorf("error login happend: %w", err)
			}
			dataId := ctx.Args().Get(2)
			data, err := store.GetData(dataId, id)
			if err != nil {
				return fmt.Errorf("error get happend: %w", err)
			}
			fmt.Println("DataID: " + data.DataID + " Data: " + data.Data + " Meta Info: " + data.Metadata)
			return nil
		},
	}
}

// DelData - used to delete data
func DelData(store storage.Storage) *cli.Command {
	return &cli.Command{
		Name:    "Delete data",
		Usage:   "used to delete data; you need to enter login and password, then data name; example: go run main.go del login password dataId",
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
			id, err := store.Login(login, password)
			if err != nil {
				return fmt.Errorf("error login happend: %w", err)
			}
			dataId := ctx.Args().Get(2)
			err = store.DelData(dataId, id)
			if err != nil {
				return fmt.Errorf("error deletr happend: %w", err)
			}
			fmt.Println("data deleted successfully")
			return nil
		},
	}
}

// Sync - used synchronize server and client
func Sync(store storage.Storage) *cli.Command {
	return &cli.Command{
		Name:    "synchronization",
		Usage:   "used synchronize server and client; you need to enter login and password; example: go run main.go sync login password",
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
			id, err := store.Login(login, password)
			if err != nil {
				return fmt.Errorf("error login happend: %w", err)
			}
			err = store.ClientSync(id, nil)
			if err != nil {
				return fmt.Errorf("error client sync happend: %w", err)
			}
			data, err := store.Sync(id)
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

// MainAction - shows help by default when app started
func MainAction(ctx *cli.Context) error {
	ctx.App.Command("help").Run(ctx)
	return nil
}

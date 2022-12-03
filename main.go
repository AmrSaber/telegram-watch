package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/AmrSaber/tw/internal/commands"
	"github.com/AmrSaber/tw/internal/common"
	"github.com/urfave/cli/v2"
)

//go:embed .env
var env string

func main() {
	common.SetEvn(env)

	app := &cli.App{
		Name:  "tw",
		Usage: "(Telegram Watch) Watch over provided command and send results to telegram",

		Commands: []*cli.Command{
			{
				Name:        "register",
				Usage:       "Registers user's info, overwriting old info if any",
				Description: "Saves user info in ~/.config/tw.yaml all of user info is saved plain except telegram id which is saved encrypted",

				Action: commands.RegisterCommand,
			},

			{
				Name:  "status",
				Usage: "Get registered user's info",

				Action: func(c *cli.Context) error {
					return nil
				},
			},

			{
				Name:  "watch",
				Usage: "Run provided command and watch it through telegram",

				Action: func(c *cli.Context) error {
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

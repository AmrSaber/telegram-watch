package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/AmrSaber/tw/internal/commands"
	"github.com/AmrSaber/tw/internal/env"
	"github.com/urfave/cli/v2"
)

//go:embed .env
var envContent string

func main() {
	env.SetEvn(envContent)

	app := &cli.App{
		Name:  "tw",
		Usage: "(Telegram Watch) Watch over provided command and send results to telegram",

		Commands: []*cli.Command{
			{
				Name:  "user",
				Usage: "Manage saved user",

				Subcommands: []*cli.Command{
					{
						Name:        "register",
						Usage:       "Registers user's info, overwriting old info if any",
						Description: "Saves user info in {user-config-directory}/tw.yaml; all of user info is saved as plain text except telegram id which is encrypted",

						Action: commands.RegisterUserCommand,
					},

					{
						Name:  "info",
						Usage: "Get registered user's info",

						Action: commands.UserInfoCommand,
					},

					{
						Name:  "delete",
						Usage: "Delete registered user's info",

						Action: commands.DeleteUserCommand,
					},
				},
			},

			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "Run provided command and pipe its output it to telegram",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "quiet",
						Aliases: []string{"q"},

						Usage: "if provided, will not show output in terminal",

						Value:    false,
						Required: false,
					},

					&cli.StringFlag{
						Name:    "timeout",
						Aliases: []string{"t"},

						Usage: "if provided, will set timeout for the running command, acceptable suffixes are (ns, ms, s, m, h) e.g. 2s, 100ms, ...",

						Value:    "",
						Required: false,
					},
				},

				Action: commands.RunCommand,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

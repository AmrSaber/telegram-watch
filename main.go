package main

import (
	_ "embed"
	"os"

	"github.com/AmrSaber/tw/internal/commands"
	"github.com/AmrSaber/tw/internal/env"
	"github.com/fatih/color"
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
						Usage:       "Registers or updates user's info",
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
				Name:        "send",
				Aliases:     []string{"notify", "message"},
				Usage:       "Sends a message to telegram",
				Description: "Can be used to send done message when a task is done or error if there is an error, can use - to send stdin",

				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "fail",
						Aliases: []string{"error", "err", "boo"},
						Usage:   "Prefix the message with a red x",
					},

					&cli.BoolFlag{
						Name:    "success",
						Aliases: []string{"done", "yay"},
						Usage:   "Prefix the message with a green check",
					},
				},

				Action: commands.SendCommand,
			},

			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "Run provided command and pipe its output to telegram",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "quiet",
						Aliases: []string{"q", "silent", "s"},

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
		color.Red("error: %s", err)
	}
}

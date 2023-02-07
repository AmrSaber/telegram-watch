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
				Name:    "notify",
				Aliases: []string{"n"},
				Usage:   "Send a message to telegram, can be used to send done message when a task is done or error if there is an error",

				Subcommands: []*cli.Command{
					{
						Name:    "done",
						Aliases: []string{"d", "success"},
						Usage:   "Sends success message to telegram, message from arguments, defaults to 'Done'",

						Action: commands.NotifyDone,
					},

					{
						Name:    "error",
						Aliases: []string{"err", "e"},
						Usage:   "Sends error message to telegram, message from arguments, defaults to 'Error'",

						Action: commands.NotifyErr,
					},

					{
						Name:    "message",
						Aliases: []string{"msg", "m"},
						Usage:   "sends custom message to telegram; the message is the argument provided to this command",

						Action: commands.NotifyMessage,
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
		color.Red("error: %s", err)
	}
}

package commands

import (
	"fmt"
	"os"
	"os/user"

	"github.com/AmrSaber/tw/internal/models"
	"github.com/AmrSaber/tw/internal/ui"
	"github.com/gosuri/uilive"
	"github.com/urfave/cli/v2"
)

func RegisterCommand(c *cli.Context) error {
	uiWriter := uilive.New()
	uiWriter.Start()

	enterInfoMsg := "Please enter your info"

	config := models.LoadConfig()

	for {
		fmt.Fprintln(uiWriter, enterInfoMsg)
		uiWriter.Flush()

		defaultName := ""
		if config.User != nil {
			defaultName = config.User.Name
		} else if currentUser, err := user.Current(); err == nil {
			defaultName = currentUser.Name
		}
		name := ui.AskString("Name:", defaultName)

		fmt.Fprintln(uiWriter, enterInfoMsg)
		fmt.Fprintln(uiWriter.Newline(), "Name:", name)
		uiWriter.Flush()

		defaultHostname, _ := os.Hostname()
		if config.User != nil {
			defaultHostname = config.User.Hostname
		}
		hostname := ui.AskString("Hostname:", defaultHostname)

		fmt.Fprintln(uiWriter, enterInfoMsg)
		fmt.Fprintln(uiWriter.Newline(), "Name:", name)
		fmt.Fprintln(uiWriter.Newline(), "Hostname:", hostname)
		uiWriter.Flush()

		defaultTelegramId := ""
		if config.User != nil {
			if id, err := config.User.DecryptTelegramId(); err == nil {
				defaultTelegramId = id
			}
		}
		telegramId := ui.AskString("Telegram ID:", defaultTelegramId)

		displayedTelegramId := ""
		for range telegramId {
			displayedTelegramId += "*"
		}

		fmt.Fprintln(uiWriter, enterInfoMsg)
		fmt.Fprintln(uiWriter.Newline(), "Name:", name)
		fmt.Fprintln(uiWriter.Newline(), "Hostname:", hostname)
		fmt.Fprintln(uiWriter.Newline(), "Telegram ID:", displayedTelegramId)
		uiWriter.Flush()

		if ui.AskBool("Save this data?", true) {
			if config.User == nil {
				config.User = &models.User{}
			}

			config.User.Name = name
			config.User.Hostname = hostname
			err := config.User.SetTelegramId(telegramId)
			if err != nil {
				fmt.Fprintln(uiWriter, "Error saving user ✗")
				uiWriter.Flush()
				uiWriter.Stop()
				return err
			}

			break
		}
	}

	if err := config.Save(); err != nil {
		fmt.Fprintln(uiWriter, "Error saving user ✗")
		uiWriter.Flush()
		uiWriter.Stop()
		return err
	}

	fmt.Fprintln(uiWriter, "User registered ✔")
	uiWriter.Flush()
	uiWriter.Stop()

	return nil
}

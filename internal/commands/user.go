package commands

import (
	"fmt"
	"os"
	"os/user"

	"github.com/AmrSaber/tw/internal/models"
	"github.com/AmrSaber/tw/internal/ui"
	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"github.com/urfave/cli/v2"
)

func RegisterUserCommand(c *cli.Context) error {
	uiWriter := uilive.New()
	uiWriter.Start()

	failUiWriter := func() {
		fmt.Fprintln(uiWriter, color.RedString("Error saving user ✘"))
		uiWriter.Flush()
		uiWriter.Stop()
	}

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
		name := ui.AskString(ui.NamePrompt, defaultName)

		fmt.Fprintln(uiWriter, enterInfoMsg)
		fmt.Fprintln(uiWriter.Newline(), ui.NamePrompt, name)
		uiWriter.Flush()

		defaultHostname, _ := os.Hostname()
		if config.User != nil {
			defaultHostname = config.User.Hostname
		}
		hostname := ui.AskString(ui.HostnamePrompt, defaultHostname)

		fmt.Fprintln(uiWriter, enterInfoMsg)
		fmt.Fprintln(uiWriter.Newline(), ui.NamePrompt, name)
		fmt.Fprintln(uiWriter.Newline(), ui.HostnamePrompt, hostname)
		uiWriter.Flush()

		defaultTelegramId := ""
		if config.User != nil {
			if id, err := config.User.DecryptTelegramId(); err == nil {
				defaultTelegramId = id
			}
		}
		telegramId := ui.AskString(ui.TelegramIdPrompt, defaultTelegramId)

		displayedTelegramId := ""
		for range telegramId {
			displayedTelegramId += "*"
		}

		fmt.Fprintln(uiWriter, enterInfoMsg)
		fmt.Fprintln(uiWriter.Newline(), ui.NamePrompt, name)
		fmt.Fprintln(uiWriter.Newline(), ui.HostnamePrompt, hostname)
		fmt.Fprintln(uiWriter.Newline(), ui.TelegramIdPrompt, displayedTelegramId)
		uiWriter.Flush()

		if ui.AskBool("Save this data?", true) {
			if config.User == nil {
				config.User = &models.User{}
			}

			config.User.Name = name
			config.User.Hostname = hostname
			err := config.User.SetTelegramId(telegramId)
			if err != nil {
				failUiWriter()
				return err
			}

			break
		}
	}

	if err := config.Save(); err != nil {
		failUiWriter()
		return err
	}

	fmt.Fprintln(uiWriter, color.GreenString("User registered ✔"))
	uiWriter.Flush()
	uiWriter.Stop()

	return nil
}

func UserInfoCommand(c *cli.Context) error {
	config := models.LoadConfig()
	if config.User == nil {
		fmt.Println(color.RedString("No saved user ✘"))
		return nil
	}

	telegramId, err := config.User.DecryptTelegramId()
	if config.User == nil {
		fmt.Println(color.RedString("Error decrypting telegram ID"))
		return err
	}

	fmt.Println(ui.NamePrompt, config.User.Name)
	fmt.Println(ui.HostnamePrompt, config.User.Hostname)
	fmt.Println(ui.TelegramIdPrompt, telegramId)

	return nil
}

func DeleteUserCommand(c *cli.Context) error {
	config := models.LoadConfig()
	if config.User == nil {
		fmt.Println(color.RedString("No saved user ✘"))
		return nil
	}

	config.User = nil
	if err := config.Save(); err != nil {
		return err
	}

	fmt.Println(color.GreenString("User info deleted ✔"))

	return nil
}

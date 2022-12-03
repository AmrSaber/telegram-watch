package commands

import (
	"fmt"

	"github.com/AmrSaber/tw/internal/models"
	"github.com/AmrSaber/tw/internal/ui"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func InfoCommand(c *cli.Context) error {
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

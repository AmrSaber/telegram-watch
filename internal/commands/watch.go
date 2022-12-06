package commands

import (
	"fmt"
	"strings"

	"github.com/AmrSaber/tw/internal/models"
	"github.com/urfave/cli/v2"
)

func WatchCommand(c *cli.Context) error {
	config := models.LoadConfig()
	if config.User == nil {
		return fmt.Errorf("no registered user")
	}

	telegramId, err := config.User.DecryptTelegramId()
	if err != nil || telegramId == "" {
		return fmt.Errorf("no registered telegram id")
	}

	args := c.Args().Slice()
	if len(args) == 0 {
		return fmt.Errorf("you must provide a command to run")
	}

	command := strings.Join(args, " ")

	watcher, err := models.NewWatcher(config, command)
	if err != nil {
		return err
	}

	if err := watcher.RunCommand(); err != nil {
		return err
	}

	return nil
}

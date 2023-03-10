package commands

import (
	"fmt"
	"os"
	"strings"
	"time"

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

	intervalStr := c.String("interval")
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		return fmt.Errorf("invalid interval: %w", err)
	}

	if interval < 4*time.Second {
		fmt.Fprintln(os.Stderr, "too low interval, setting it to 4 seconds!")
		interval = 4 * time.Second
	}

	timeoutStr := c.String("timeout")
	var timeout time.Duration
	if timeoutStr != "" {
		var err error
		timeout, err = time.ParseDuration(timeoutStr)
		if err != nil {
			return fmt.Errorf("invalid timeout: %w", err)
		}
	}

	quiet := c.Bool("quiet")

	config.Runtime.Quiet = quiet
	config.Runtime.Interval = interval
	config.Runtime.Timeout = timeout

	args := c.Args().Slice()
	if len(args) == 0 {
		return fmt.Errorf("you must provide a command to run")
	}

	command := strings.Join(args, " ")
	watcher := models.NewWatcher(config, command)
	if err := watcher.WatchCommand(); err != nil {
		return err
	}

	fmt.Println("\nDone watching", command)

	return nil
}

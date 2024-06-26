package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AmrSaber/tw/internal/models"
	"github.com/AmrSaber/tw/internal/utils"
	"github.com/urfave/cli/v2"
)

func RunCommand(c *cli.Context) error {
	quiet := c.Bool("quiet")
	timeout := c.String("timeout")

	config := models.LoadConfig()
	if config.User == nil {
		return fmt.Errorf("no registered user")
	}

	config.Runtime.Quiet = quiet

	telegramId, err := config.User.DecryptTelegramId()
	if err != nil || telegramId == "" {
		return fmt.Errorf("no registered telegram id")
	}

	args := c.Args().Slice()
	if len(args) == 0 {
		return fmt.Errorf("you must provide a command to run")
	}

	command := strings.Join(args, " ")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if timeout != "" {
		duration, err := time.ParseDuration(timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout: %w", err)
		}

		go func() {
			time.Sleep(duration)
			cancel()
		}()
	}

	utils.HandleInterrupt(func() { cancel() })

	runner, err := models.NewRunner(config, command)
	if err != nil {
		return err
	}

	return runner.RunCommand(ctx)
}

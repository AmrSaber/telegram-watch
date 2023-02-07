package commands

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/AmrSaber/tw/internal/models"
	"github.com/AmrSaber/tw/internal/utils"
	"github.com/urfave/cli/v2"
)

func SendCommand(ctx *cli.Context) error {
	if ctx.Args().Len() == 0 {
		return fmt.Errorf("you must provide a message")
	}

	// Read message (from args or std)
	message := strings.Join(ctx.Args().Slice(), " ")
	message = strings.TrimSpace(message)
	if message == "-" {
		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		message = string(bytes)
	}

	// Add any required prefixes
	if ctx.Bool("done") {
		message = fmt.Sprintf("%s %s", utils.GREEN_CHECK, message)
	} else if ctx.Bool("err") {
		message = fmt.Sprintf("%s %s", utils.RED_X, message)
	}

	// Load and validate config
	config := models.LoadConfig()
	if config.User == nil {
		return fmt.Errorf("no registered user")
	}

	telegramId, err := config.User.DecryptTelegramId()
	if err != nil || telegramId == "" {
		return fmt.Errorf("no registered telegram id")
	}

	// Send message
	message = fmt.Sprintf("%s:\n%s", config.User.Hostname, message)
	return utils.SendSingleMessage(telegramId, message)
}

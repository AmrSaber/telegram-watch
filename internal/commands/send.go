package commands

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unsafe"

	"github.com/AmrSaber/tw/internal/models"
	"github.com/AmrSaber/tw/internal/utils"
	"github.com/urfave/cli/v2"
)

func SendCommand(ctx *cli.Context) error {
	if ctx.Args().Len() == 0 {
		return fmt.Errorf("you must provide a message")
	}

	// Read message content (from args or std)
	content := strings.Join(ctx.Args().Slice(), " ")
	content = strings.TrimSpace(content)
	if content == "-" {
		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		content = string(bytes)
	}

	// Add any required prefixes
	if ctx.Bool("done") {
		content = fmt.Sprintf("%s %s", utils.GREEN_CHECK, content)
	} else if ctx.Bool("err") {
		content = fmt.Sprintf("%s %s", utils.RED_X, content)
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
	content = fmt.Sprintf("%s:\n%s", config.User.Hostname, content)

	message := models.NewTelegramMessage(telegramId)
	_, err = message.Write(unsafe.Slice(unsafe.StringData(content), len(content)))

	return err
}

package commands

import (
	"fmt"
	"strings"

	"github.com/AmrSaber/tw/internal/models"
	"github.com/AmrSaber/tw/internal/utils"
	"github.com/urfave/cli/v2"
)

func NotifyDone(ctx *cli.Context) error {
	message := fmt.Sprintf("%s Done", utils.GREEN_CHECK)
	if ctx.Args().Len() != 0 {
		args := strings.Join(ctx.Args().Slice(), " ")
		message = fmt.Sprintf("%s %s", utils.GREEN_CHECK, args)
	}

	return sendMessage(message)
}

func NotifyErr(ctx *cli.Context) error {
	message := fmt.Sprintf("%s Error", utils.RED_X)
	if ctx.Args().Len() != 0 {
		args := strings.Join(ctx.Args().Slice(), " ")
		message = fmt.Sprintf("%s %s", utils.RED_X, args)
	}

	return sendMessage(message)
}

func NotifyMessage(ctx *cli.Context) error {
	if ctx.Args().Len() == 0 {
		return fmt.Errorf("you must provide a message")
	}

	args := strings.Join(ctx.Args().Slice(), " ")
	return sendMessage(args)
}

func sendMessage(message string) error {
	config := models.LoadConfig()
	if config.User == nil {
		return fmt.Errorf("no registered user")
	}

	telegramId, err := config.User.DecryptTelegramId()
	if err != nil || telegramId == "" {
		return fmt.Errorf("no registered telegram id")
	}

	message = fmt.Sprintf("%s:\n%s", config.User.Hostname, message)

	return utils.SendSingleMessage(telegramId, message)
}

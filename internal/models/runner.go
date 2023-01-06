package models

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/AmrSaber/tw/internal/env"
	"github.com/AmrSaber/tw/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandRunner struct {
	bot          *tgbotapi.BotAPI
	stdoutWriter *TelegramWriter
	stderrWriter *TelegramWriter

	command *exec.Cmd

	doneSuccessMessage string
	doneFailMessage    string
}

var messageBaseTemplate string

func NewRunner(ctx context.Context, config Config, command string) (*CommandRunner, error) {
	cmd := exec.CommandContext(ctx, "bash", "-c", command)

	bot, err := tgbotapi.NewBotAPI(env.GetBotTokenKey())
	if err != nil {
		return nil, err
	}

	messageBaseTemplate = fmt.Sprintf(
		"%%s Hello %s\n"+
			"This message traces %%s from command %q\n"+
			"From your device %q\n\n"+
			"%%%%s",
		config.User.Name,
		command,
		config.User.Hostname,
	)

	stdoutTemplate := fmt.Sprintf(messageBaseTemplate, utils.GREEN_CIRCLE, "STDOUT")
	stderrTemplate := fmt.Sprintf(messageBaseTemplate, utils.RED_CIRCLE, "STDERR")

	telegramId, _ := config.User.DecryptTelegramId()
	stdoutMsgWriter := NewTelegramWriter(telegramId, bot, stdoutTemplate)
	stderrMsgWriter := NewTelegramWriter(telegramId, bot, stderrTemplate)

	var stdoutWriter, stderrWriter io.Writer
	stdoutWriter = stdoutMsgWriter
	stderrWriter = stderrMsgWriter

	if !config.Runtime.Quiet {
		stdoutWriter = io.MultiWriter(os.Stdout, stdoutMsgWriter)
		stderrWriter = io.MultiWriter(os.Stderr, stderrMsgWriter)
	}

	cmd.Stdout = stdoutWriter
	cmd.Stderr = stderrWriter

	doneSuccessMessage := fmt.Sprintf(
		"Command %q\non device %q\nis now done successfully %s",
		command,
		config.User.Hostname,
		utils.GREEN_CHECK,
	)

	doneFailMessage := fmt.Sprintf(
		"Command %q\non device %q\nexited with %%q %s",
		command,
		config.User.Hostname,
		utils.RED_X,
	)

	runner := CommandRunner{
		bot:          bot,
		stdoutWriter: stdoutMsgWriter,
		stderrWriter: stderrMsgWriter,

		command: cmd,

		doneSuccessMessage: doneSuccessMessage,
		doneFailMessage:    doneFailMessage,
	}

	return &runner, nil
}

func (w *CommandRunner) RunCommand() error {
	err := w.command.Run()

	// Set templates to completed templates
	stdoutTemplate := fmt.Sprintf(messageBaseTemplate, utils.WHITE_CIRCLE, "STDOUT")
	stderrTemplate := fmt.Sprintf(messageBaseTemplate, utils.WHITE_CIRCLE, "STDERR")
	w.stdoutWriter.SetTemplate(stdoutTemplate)
	w.stderrWriter.SetTemplate(stderrTemplate)

	doneMsgConfig := tgbotapi.NewMessage(w.stdoutWriter.chatId, w.doneSuccessMessage)
	if err != nil {
		failMessage := fmt.Sprintf(w.doneFailMessage, err)
		doneMsgConfig = tgbotapi.NewMessage(w.stderrWriter.chatId, failMessage)
	}

	if _, err := w.bot.Send(doneMsgConfig); err != nil {
		return err
	}

	return nil
}

package models

import (
	"fmt"
	"os/exec"

	"github.com/AmrSaber/tw/internal/env"
	"github.com/AmrSaber/tw/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandWatcher struct {
	bot          *tgbotapi.BotAPI
	stdoutWriter *TelegramWriter
	stderrWriter *TelegramWriter

	command *exec.Cmd

	doneSuccessMessage string
	doneFailMessage    string
}

var messageBaseTemplate string

func NewWatcher(config Config, command string) (*CommandWatcher, error) {
	cmd := exec.Command("bash", "-c", command)

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

	stdoutTemplate := fmt.Sprintf(messageBaseTemplate, utils.BLUE_CIRCLE, "STDOUT")
	stderrTemplate := fmt.Sprintf(messageBaseTemplate, utils.RED_CIRCLE, "STDERR")

	telegramId, _ := config.User.DecryptTelegramId()
	stdoutWriter := NewTelegramWriter(telegramId, bot, stdoutTemplate)
	stderrWriter := NewTelegramWriter(telegramId, bot, stderrTemplate)

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

	watcher := CommandWatcher{
		bot:          bot,
		stdoutWriter: stdoutWriter,
		stderrWriter: stderrWriter,

		command: cmd,

		doneSuccessMessage: doneSuccessMessage,
		doneFailMessage:    doneFailMessage,
	}

	return &watcher, nil
}

func (w *CommandWatcher) RunCommand() error {
	err := w.command.Run()

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

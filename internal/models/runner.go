package models

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

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

func NewRunner(config Config, command string) (*CommandRunner, error) {
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

	stdoutTemplate := fmt.Sprintf(messageBaseTemplate, utils.GREEN_CIRCLE, "STDOUT")
	stderrTemplate := fmt.Sprintf(messageBaseTemplate, utils.RED_CIRCLE, "STDERR")

	telegramId, _ := config.User.DecryptTelegramId()

	stdoutMsgWriter := NewTelegramWriter(telegramId)
	stdoutMsgWriter.SetContentMapper(func(input []byte) []byte { return fmt.Appendf([]byte{}, stdoutTemplate, input) })

	stderrMsgWriter := NewTelegramWriter(telegramId)
	stderrMsgWriter.SetContentMapper(func(input []byte) []byte { return fmt.Appendf([]byte{}, stderrTemplate, input) })

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

func (r *CommandRunner) RunCommand(ctx context.Context) error {
	err := r.command.Start()
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		r.command.Process.Signal(syscall.SIGINT)
	}()

	err = r.command.Wait()

	// Wait for writers to finish any pending writing
	r.stdoutWriter.Wait()
	r.stderrWriter.Wait()

	// Set templates to completed templates
	stdoutTemplate := fmt.Sprintf(messageBaseTemplate, utils.WHITE_CIRCLE, "STDOUT")
	r.stdoutWriter.SetContentMapper(func(input []byte) []byte { return fmt.Appendf([]byte{}, stdoutTemplate, input) })

	stderrTemplate := fmt.Sprintf(messageBaseTemplate, utils.WHITE_CIRCLE, "STDERR")
	r.stderrWriter.SetContentMapper(func(input []byte) []byte { return fmt.Appendf([]byte{}, stderrTemplate, input) })

	doneMsgConfig := tgbotapi.NewMessage(r.stdoutWriter.GetChatId(), r.doneSuccessMessage)
	if err != nil {
		failMessage := fmt.Sprintf(r.doneFailMessage, err)
		doneMsgConfig = tgbotapi.NewMessage(r.stderrWriter.GetChatId(), failMessage)
	}

	if _, err := r.bot.Send(doneMsgConfig); err != nil {
		return err
	}

	return err
}

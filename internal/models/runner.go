package models

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/AmrSaber/tw/internal/utils"
)

type CommandRunner struct {
	stdoutWriter *TelegramWriter
	stderrWriter *TelegramWriter

	command *exec.Cmd

	doneSuccessMessage string
	doneFailMessage    string
}

var runnerMessageBaseTemplate string

func NewRunner(config Config, command string) (*CommandRunner, error) {
	cmd := exec.Command("bash", "-c", command)

	runnerMessageBaseTemplate = fmt.Sprintf(
		"%%s Hello %s\n"+
			"This message traces %%s from command %q\n"+
			"From your device %q\n\n"+
			"%%%%s",
		config.User.Name,
		command,
		config.User.Hostname,
	)

	stdoutTemplate := fmt.Sprintf(runnerMessageBaseTemplate, utils.GREEN_CIRCLE, "STDOUT")
	stderrTemplate := fmt.Sprintf(runnerMessageBaseTemplate, utils.RED_CIRCLE, "STDERR")

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
		"Command %q\non device %q\nis now done successfully %s\nTotal run time %%s",
		command,
		config.User.Hostname,
		utils.GREEN_CHECK,
	)

	doneFailMessage := fmt.Sprintf(
		"Command %q\non device %q\nexited with %%q %s\nTotal run time %%s",
		command,
		config.User.Hostname,
		utils.RED_X,
	)

	runner := CommandRunner{
		stdoutWriter: stdoutMsgWriter,
		stderrWriter: stderrMsgWriter,

		command: cmd,

		doneSuccessMessage: doneSuccessMessage,
		doneFailMessage:    doneFailMessage,
	}

	return &runner, nil
}

func (r *CommandRunner) RunCommand(ctx context.Context) error {
	commandStartTime := time.Now()

	err := r.command.Start()
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		r.command.Process.Signal(syscall.SIGINT)
	}()

	err = r.command.Wait()

	commandRunTime := time.Since(commandStartTime)

	// Set templates to completed templates
	stdoutTemplate := fmt.Sprintf(runnerMessageBaseTemplate, utils.WHITE_CIRCLE, "STDOUT")
	r.stdoutWriter.SetContentMapper(func(input []byte) []byte { return fmt.Appendf([]byte{}, stdoutTemplate, input) })

	stderrTemplate := fmt.Sprintf(runnerMessageBaseTemplate, utils.WHITE_CIRCLE, "STDERR")
	r.stderrWriter.SetContentMapper(func(input []byte) []byte { return fmt.Appendf([]byte{}, stderrTemplate, input) })

	doneMessage := fmt.Sprintf(r.doneSuccessMessage, utils.FormatTime(commandRunTime.Nanoseconds()))
	if err != nil {
		doneMessage = fmt.Sprintf(r.doneFailMessage, err, utils.FormatTime(commandRunTime.Nanoseconds()))
	}

	// Wait for writers to finish any pending writing
	r.stdoutWriter.Wait()
	r.stderrWriter.Wait()

	doneWriter := NewTelegramWriter(r.stdoutWriter.GetChatId())
	doneWriter.Write(utils.StringToBytes(doneMessage))
	doneWriter.Wait()

	return err
}

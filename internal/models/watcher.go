package models

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/AmrSaber/tw/internal/utils"
	"github.com/gosuri/uilive"
)

type CommandWatcher struct {
	messageWriter *TelegramWriter

	command       string
	doneMessage   string
	runtimeConfig *RuntimeConfig
}

var watcherMessageBaseTemplate string

func NewWatcher(config Config, command string) *CommandWatcher {

	watcherMessageBaseTemplate = fmt.Sprintf(
		"%%s Hello %s\n"+
			"This message traces watching output from command %q\n"+
			"From your device %q\n"+
			"Last update at %%%%s\n\n"+
			"%%%%s",
		config.User.Name,
		command,
		config.User.Hostname,
	)

	outputTemplate := fmt.Sprintf(watcherMessageBaseTemplate, utils.GREEN_CIRCLE)

	telegramId, _ := config.User.DecryptTelegramId()

	messageWriter := NewTelegramWriter(telegramId)
	messageWriter.SetContentMapper(func(input []byte) []byte {
		now := time.Now()
		return fmt.Appendf([]byte{}, outputTemplate, now.Format(time.RFC3339), input)
	})

	doneMessage := fmt.Sprintf(
		"%s Done watching command %q\non device %q",
		utils.GREEN_CHECK,
		command,
		config.User.Hostname,
	)

	return &CommandWatcher{
		messageWriter: messageWriter,

		command:     command,
		doneMessage: doneMessage,

		runtimeConfig: &config.Runtime,
	}
}

func (r *CommandWatcher) WatchCommand() error {
	var runningProcess *os.Process
	looping := true

	utils.HandleInterrupt(func() {
		if looping {
			looping = false

			// stop the running process
			if runningProcess != nil {
				runningProcess.Signal(syscall.SIGINT)
				runningProcess = nil
			}
		} else {
			os.Exit(1)
		}
	})

	var uiWriter *uilive.Writer
	if !r.runtimeConfig.Quiet {
		uiWriter = uilive.New()
		defer uiWriter.Stop()
	}

	var outputBuffer bytes.Buffer

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var timer *time.Timer
	if r.runtimeConfig.Timeout != 0 {
		timer = time.NewTimer(r.runtimeConfig.Timeout)
		defer timer.Stop()

		go func() {
			for {
				select {
				case <-timer.C:
					if runningProcess != nil {
						runningProcess.Signal(syscall.SIGINT)
						runningProcess = nil
					}
					timer.Stop()
					timer.Reset(r.runtimeConfig.Timeout)

				case <-ctx.Done():
					return
				}
			}
		}()
	}

	fmt.Printf("Watching command %q\n\n", r.command)
	if !r.runtimeConfig.Quiet {
		uiWriter.Start()
	}

	for looping {
		cmd := exec.Command("bash", "-c", r.command)

		cmd.Stdout = &outputBuffer
		cmd.Stderr = &outputBuffer

		commandStartTime := time.Now()
		err := cmd.Start()
		if err != nil {
			return fmt.Errorf("could not run command: %w", err)
		}

		runningProcess = cmd.Process

		err = cmd.Wait()
		if err != nil {
			return fmt.Errorf("error waiting for command: %w", err)
		}

		if !looping {
			break
		}

		// Append run time to the output
		commandRunTime := time.Since(commandStartTime)
		outputBytes := bytes.Clone(bytes.TrimSpace(outputBuffer.Bytes()))
		outputBytes = fmt.Appendf(outputBytes, "\n\n--------------\nRun time: %s\n", utils.FormatTime(commandRunTime.Nanoseconds()))

		if !r.runtimeConfig.Quiet {
			uiWriter.Write(outputBytes)
			uiWriter.Flush()
		}

		r.messageWriter.SetContent(outputBytes)
		outputBuffer.Reset()

		time.Sleep(r.runtimeConfig.Interval)

		// Reset timer
		if timer != nil {
			if !timer.Stop() {
				<-timer.C
			}

			timer.Reset(r.runtimeConfig.Timeout)
		}
	}

	// Set templates to completed templates
	stdoutTemplate := fmt.Sprintf(watcherMessageBaseTemplate, utils.WHITE_CIRCLE)
	r.messageWriter.SetContentMapper(func(input []byte) []byte {
		now := time.Now()
		return fmt.Appendf([]byte{}, stdoutTemplate, now.Format(time.RFC3339), input)
	})

	// Wait for writers to finish any pending writing
	r.messageWriter.Wait()

	doneMessage, chatId := r.doneMessage, fmt.Sprint(r.messageWriter.GetChatId())
	doneWriter := NewTelegramWriter(chatId)
	doneWriter.Write(utils.ToBytes(doneMessage))
	doneWriter.Wait()

	fmt.Printf("\nDone watching %q\n", r.command)

	return nil
}

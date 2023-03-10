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
		} else if runningProcess != nil {
			runningProcess.Signal(syscall.SIGINT)
			runningProcess = nil
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

	timer := time.NewTimer(r.runtimeConfig.Timeout)
	defer timer.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case <-timer.C:
				fmt.Println("timer timed out")
				fmt.Println("running process:", runningProcess)
				// FIXME runningProcess is always nil
				if runningProcess != nil {
					fmt.Println("signaled")
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

	fmt.Printf("Watching command %s\n\n", r.command)
	uiWriter.Start()

	for looping {
		cmd := exec.Command("bash", "-c", r.command)
		runningProcess = cmd.Process

		cmd.Stdout = &outputBuffer
		cmd.Stderr = &outputBuffer

		startTime := time.Now()
		err := cmd.Start()
		if err != nil {
			return fmt.Errorf("could not run command: %w", err)
		}

		err = cmd.Wait()
		if err != nil {
			return fmt.Errorf("error waiting for command: %w", err)
		}

		// Append run time to the output
		commandRunTime := time.Since(startTime)
		fmt.Fprintf(&outputBuffer, "\n\n--------------\nRun time: %.3fms", float64(commandRunTime.Microseconds())/1000)

		if !r.runtimeConfig.Quiet {
			uiWriter.Write(outputBuffer.Bytes())
			uiWriter.Flush()
		}

		r.messageWriter.SetContent(outputBuffer.Bytes())
		outputBuffer.Reset()

		time.Sleep(r.runtimeConfig.Interval)

		// Reset timer
		if !timer.Stop() {
			<-timer.C
		}
		timer.Reset(r.runtimeConfig.Timeout)
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
	NewTelegramWriter(chatId).Write(utils.ToBytes(doneMessage))

	return nil
}

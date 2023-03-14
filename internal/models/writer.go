package models

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/AmrSaber/tw/internal/env"
	"github.com/AmrSaber/tw/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramWriter struct {
	lock *sync.Mutex

	chatId int64

	cooldown    bool
	shouldFlush bool
	flushed     *sync.Cond

	bot      *tgbotapi.BotAPI
	messages []*tgbotapi.Message

	fullContent   []byte
	contentMapper func([]byte) []byte
}

func NewTelegramWriter(chatId string) *TelegramWriter {
	bot, err := tgbotapi.NewBotAPI(env.GetBotTokenKey())
	if err != nil {
		panic(err)
	}

	numChatId, err := strconv.ParseInt(chatId, 0, 0)
	if err != nil {
		panic(err)
	}

	lock := sync.Mutex{}

	return &TelegramWriter{
		lock: &lock,

		chatId: numChatId,

		bot:      bot,
		messages: make([]*tgbotapi.Message, 0),

		fullContent: make([]byte, 0),

		flushed: sync.NewCond(&lock),
	}
}

func (w *TelegramWriter) GetChatId() string {
	return fmt.Sprint(w.chatId)
}

// Waits for pending message writes
func (w *TelegramWriter) Wait() {
	w.flushed.L.Lock()
	defer w.flushed.L.Unlock()

	for w.shouldFlush {
		w.flushed.Wait()
	}
}

func (w *TelegramWriter) SetCooldown(duration time.Duration) {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.cooldown = true

	go func() {
		time.Sleep(duration)

		w.lock.Lock()
		defer w.lock.Unlock()

		w.cooldown = false
		if w.shouldFlush {
			w.flush()
		}
	}()
}

func (w *TelegramWriter) SetContentMapper(mapper func([]byte) []byte) {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.contentMapper = mapper
	w.flush()
}

// Never returns an error
func (w *TelegramWriter) Write(input []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.fullContent = append(w.fullContent, input...)
	w.flush()

	return len(input), nil
}

func (w *TelegramWriter) SetContent(content []byte) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if bytes.Equal(w.fullContent, content) {
		return
	}

	w.fullContent = bytes.Clone(content)
	w.flush()
}

// Flush should always be called with lock
func (w *TelegramWriter) flush() {
	if len(w.fullContent) == 0 {
		return
	}

	w.shouldFlush = true

	if w.cooldown {
		return
	}

	// Throttle calling telegram API, in a go-routine for locks
	go func() { w.SetCooldown(4 * time.Second) }()
	w.cooldown = true

	err := w.handleMessages()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error sending/updating telegram messages:", err)

		// If the error is API throttle, leave `shouldFlush` as it is, to retry again later
		if strings.Contains(strings.ToLower(err.Error()), "retry after") {
			return
		}
	}

	w.shouldFlush = false
	w.flushed.Broadcast()
}

func (w *TelegramWriter) handleMessages() error {
	content := bytes.TrimSpace(w.fullContent)
	if w.contentMapper != nil {
		content = bytes.TrimSpace(w.contentMapper(w.fullContent))
	}

	splitContent := utils.SplitMeaningfully(content, utils.TELEGRAM_MESSAGE_LIMIT)
	countTemplate := fmt.Sprintf("\n--------------\n(%%0%dd/%d)\n", utils.CountDigits(len(splitContent)), len(splitContent))

	// Update existing messages or send new ones as needed
	for index, part := range splitContent {

		// If the message consists of more than 1 part, append part counter to the end of each part
		if len(splitContent) > 1 {
			part = bytes.Clone(part)
			part = fmt.Appendf(part, countTemplate, index+1)
		}

		// Convert the byte array into string without allocation
		strPart := utils.BytesToString(part)

		if index >= len(w.messages) {
			// If message at current index does not exist, send it
			msgConfig := tgbotapi.NewMessage(w.chatId, strPart)

			if index > 0 {
				msgConfig.ReplyToMessageID = w.messages[index-1].MessageID
			}

			msg, err := w.bot.Send(msgConfig)
			if err != nil {
				return err
			}

			w.messages = append(w.messages, &msg)
		} else {
			// If message already exists, update it
			message := w.messages[index]

			if strPart == message.Text {
				continue
			}

			updateMsgConfig := tgbotapi.NewEditMessageText(w.chatId, message.MessageID, strPart)

			msg, err := w.bot.Send(updateMsgConfig)
			if err != nil && !strings.Contains(err.Error(), "message is not modified") {
				return err
			}

			w.messages[index] = &msg
		}
	}

	// Delete any extra messages after the update
	if len(w.messages) > len(splitContent) {
		extra := w.messages[len(splitContent):]
		w.messages = w.messages[:len(splitContent)]

		for _, message := range extra {
			deleteMessageConfig := tgbotapi.NewDeleteMessage(w.chatId, message.MessageID)
			if _, err := w.bot.Send(deleteMessageConfig); err != nil {
				return err
			}
		}
	}

	return nil
}

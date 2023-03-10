package models

import (
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
	lock sync.Mutex
	once sync.Once
	wg   sync.WaitGroup

	chatId int64

	shouldFlush bool

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

	return &TelegramWriter{
		chatId: numChatId,

		bot:      bot,
		messages: make([]*tgbotapi.Message, 0),

		fullContent: make([]byte, 0),
	}
}

func (w *TelegramWriter) GetChatId() int64 {
	return w.chatId
}

// Waits for pending message writes
func (w *TelegramWriter) Wait() {
	w.wg.Wait()
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

	w.fullContent = content
	w.flush()
}

func (w *TelegramWriter) flush() {
	if len(w.fullContent) == 0 {
		return
	}

	w.shouldFlush = true

	// Throttle calling telegram API
	w.once.Do(func() {
		w.wg.Add(1)
		go func() {
			defer w.wg.Done()

			time.Sleep(4 * time.Second)
			w.lock.Lock()
			defer w.lock.Unlock()

			w.once = sync.Once{}
			if w.shouldFlush {
				w.flush()
			}
		}()

		err := w.handleMessages()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error sending/updating telegram messages:", err)
		}

		if err == nil || !strings.Contains(strings.ToLower(err.Error()), "retry after") {
			w.shouldFlush = false
		}
	})
}

func (w *TelegramWriter) handleMessages() error {
	content := w.fullContent
	if w.contentMapper != nil {
		content = w.contentMapper(w.fullContent)
	}

	splitContent := utils.MeaningfullySplit(content, utils.TELEGRAM_MESSAGE_LIMIT)
	countTemplate := fmt.Sprintf("\n\n---------------\n(%%0%dd/%d)", utils.CountDigits(len(splitContent)), len(splitContent))

	// Update existing messages or send new ones as needed
	for index, part := range splitContent {

		// If the message consists of more than 1 part, append part counter to the end of each part
		if len(splitContent) > 1 {
			// This is necessary so that the original string is not modified
			{
				clone := make([]byte, len(part))
				copy(clone, part)
				part = clone
			}

			part = fmt.Appendf(part, countTemplate, index+1)
		}

		// Convert the byte array into string without allocation
		strPart := utils.ToString(part)

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
		extra := w.messages[len(splitContent)-1:]
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

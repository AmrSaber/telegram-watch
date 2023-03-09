package models

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramTemplateWriter struct {
	lock sync.Mutex

	chatId int64

	bot *tgbotapi.BotAPI
	msg *tgbotapi.Message

	content  string
	template string // Must contain exactly one "%s" to be replaces with the content
}

func NewTelegramTemplateWriter(userId string, bot *tgbotapi.BotAPI, template string) *TelegramTemplateWriter {
	chatId, err := strconv.ParseInt(userId, 0, 0)
	if err != nil {
		panic(err)
	}

	return &TelegramTemplateWriter{
		chatId: chatId,
		bot:    bot,

		template: template,
	}
}

// Sends first message to user
func (w *TelegramTemplateWriter) start() error {
	msgContent := fmt.Sprintf(w.template, w.content)
	msgConfig := tgbotapi.NewMessage(w.chatId, msgContent)

	msg, err := w.bot.Send(msgConfig)
	if err != nil {
		return err
	}

	w.msg = &msg
	return nil
}

func (w *TelegramTemplateWriter) Write(bytes []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.msg == nil {
		w.start()
	}

	w.content += string(bytes)

	if err := w.flush(); err != nil {
		return 0, err
	}

	return len(bytes), nil
}

func (w *TelegramTemplateWriter) SetTemplate(template string) {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.template = template
	if w.msg != nil {
		w.flush()
	}
}

func (w *TelegramTemplateWriter) GetChatId() int64 {
	return w.chatId
}

// Write latest content to message
func (w *TelegramTemplateWriter) flush() error {
	msgContent := fmt.Sprintf(w.template, w.content)
	updateMsgConfig := tgbotapi.NewEditMessageText(w.chatId, w.msg.MessageID, msgContent)

	msg, err := w.bot.Send(updateMsgConfig)
	if err != nil && !strings.Contains(err.Error(), "message is not modified") {
		return err
	}

	w.msg = &msg
	return nil
}

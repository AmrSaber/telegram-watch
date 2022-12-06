package models

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramWriter struct {
	chatId int64

	bot *tgbotapi.BotAPI
	msg *tgbotapi.Message

	content  string
	template string // Must contain exactly one "%s" to be replaces with the content
}

func NewTelegramWriter(userId string, bot *tgbotapi.BotAPI, template string) *TelegramWriter {
	chatId, err := strconv.ParseInt(userId, 0, 0)
	if err != nil {
		panic(err)
	}

	return &TelegramWriter{
		chatId: chatId,
		bot:    bot,

		template: template,
	}
}

// Sends first message to user
func (w *TelegramWriter) start() error {
	msgContent := fmt.Sprintf(w.template, w.content)
	msgConfig := tgbotapi.NewMessage(w.chatId, msgContent)

	msg, err := w.bot.Send(msgConfig)
	if err != nil {
		return err
	}

	w.msg = &msg
	return nil
}

func (w *TelegramWriter) Write(bytes []byte) (int, error) {
	if w.msg == nil {
		w.start()
	}

	w.content += string(bytes)

	if err := w.flush(); err != nil {
		return 0, err
	}

	return len(bytes), nil
}

func (w *TelegramWriter) SetTemplate(template string) {
	w.template = template
	if w.msg != nil {
		w.flush()
	}
}

func (w *TelegramWriter) GetChatId() int64 {
	return w.chatId
}

// Write latest content to message
func (w *TelegramWriter) flush() error {
	msgContent := fmt.Sprintf(w.template, w.content)
	updateMsgConfig := tgbotapi.NewEditMessageText(w.chatId, w.msg.MessageID, msgContent)

	msg, err := w.bot.Send(updateMsgConfig)
	if err != nil {
		return err
	}

	w.msg = &msg
	return nil
}

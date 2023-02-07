package utils

import (
	"strconv"

	"github.com/AmrSaber/tw/internal/env"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendSingleMessage(telegramId string, message string) error {
	bot, err := tgbotapi.NewBotAPI(env.GetBotTokenKey())
	if err != nil {
		return err
	}

	chatId, err := strconv.ParseInt(telegramId, 0, 0)
	if err != nil {
		return err
	}

	msgConfig := tgbotapi.NewMessage(chatId, message)

	_, err = bot.Send(msgConfig)
	return err
}

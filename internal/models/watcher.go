package models

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramWatcher struct {
	bot       *tgbotapi.BotAPI
	stdOutMsg tgbotapi.Message
	stdErrMsg tgbotapi.Message
}

func NewWatcher(config Config, command string) *TelegramWatcher {
	// bot, err := tgbotapi.NewBotAPI(env.GetBotTokenKey())
	// if err != nil {
	// 	return nil, err
	// }

	// stringId, _ := config.User.DecryptTelegramId()
	// id, _ := strconv.Atoi(stringId)
	// msgConfig := tgbotapi.NewMessage(int64(id), "Hello, this is a message to track your command")

	// msg, err := bot.Send(msgConfig)
	// if err != nil {
	// 	return nil
	// }

	return nil
}

func (w *TelegramWatcher) Close() error {
	return nil
}

func (w *TelegramWatcher) Write(bytes []byte) (int, error) {

	// w.output += string(bytes)
	// fmt.Println(string(bytes))

	// editMessageConfig := tgbotapi.NewEditMessageText(w.msg.Chat.ID, w.msg.MessageID, w.output)
	// var err error

	// w.msg, err = w.bot.Send(editMessageConfig)
	// if err != nil {
	// 	return 0, nil
	// }

	return len(bytes), nil
}

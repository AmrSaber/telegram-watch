package models

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/AmrSaber/tw/internal/env"
	"github.com/AmrSaber/tw/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramWriter struct {
	lock sync.Mutex

	chatId    int64
	autoFlush bool

	bot      *tgbotapi.BotAPI
	messages []*tgbotapi.Message

	fullContent []byte
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
		lock: sync.Mutex{},

		chatId:    numChatId,
		autoFlush: true,

		bot:      bot,
		messages: make([]*tgbotapi.Message, 0),

		fullContent: make([]byte, 0),
	}
}

func (m *TelegramWriter) GetChatId() int64 {
	return m.chatId
}

func (m *TelegramWriter) SetAutoFlush(autoFlush bool) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.autoFlush = autoFlush
}

func (m *TelegramWriter) Write(input []byte) (int, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.fullContent = append(m.fullContent, input...)

	if m.autoFlush {
		if err := m.flush(); err != nil {
			return 0, err
		}
	}

	return len(input), nil
}

func (m *TelegramWriter) SetContent(content []byte) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.fullContent = content

	if m.autoFlush {
		if err := m.flush(); err != nil {
			return err
		}
	}

	return nil
}

// FIXME: too many calls trigger rate limit (try with "ping google.com")
func (m *TelegramWriter) flush() error {
	if len(m.fullContent) == 0 {
		return nil
	}

	splitContent := utils.MeaningfullySplit(m.fullContent, utils.TELEGRAM_MESSAGE_LIMIT)
	countTemplate := fmt.Sprintf("\n\n---------------\n(%%0%dd/%d)", utils.CountDigits(len(splitContent)), len(splitContent))

	// Update existing messages or send new ones as needed
	for index, part := range splitContent {

		// If the message consists of more than 1 part, append part counter to the end of each part
		if len(splitContent) > 1 {
			// This part is necessary so that the original string is not modified
			{
				partClone := make([]byte, len(part))
				copy(partClone, part)
				part = partClone
			}

			part = fmt.Appendf(part[:], countTemplate, index+1)
		}

		// Convert the byte array into string without allocation
		strPart := utils.ToString(part)

		if index >= len(m.messages) {
			// If message at current index does not exist, send it
			msgConfig := tgbotapi.NewMessage(m.chatId, strPart)

			if index > 0 {
				msgConfig.ReplyToMessageID = m.messages[index-1].MessageID
			}

			msg, err := m.bot.Send(msgConfig)
			if err != nil {
				return err
			}

			m.messages = append(m.messages, &msg)
		} else {
			message := m.messages[index]

			if strPart == message.Text {
				continue
			}

			// If message already exists, update it
			updateMsgConfig := tgbotapi.NewEditMessageText(m.chatId, message.MessageID, strPart)

			msg, err := m.bot.Send(updateMsgConfig)
			if err != nil && !strings.Contains(err.Error(), "message is not modified") {
				return err
			}

			m.messages[index] = &msg
		}
	}

	// Delete any extra messages after the update
	if len(m.messages) > len(splitContent) {
		extra := m.messages[len(splitContent)-1:]
		m.messages = m.messages[:len(splitContent)]

		for _, message := range extra {
			deleteMessageConfig := tgbotapi.NewDeleteMessage(m.chatId, message.MessageID)
			if _, err := m.bot.Send(deleteMessageConfig); err != nil {
				return err
			}
		}
	}

	return nil
}

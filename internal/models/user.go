package models

import (
	"github.com/AmrSaber/tw/internal/env"
	"github.com/AmrSaber/tw/internal/utils"
)

type User struct {
	Name       string `yaml:"name"`
	Hostname   string `yaml:"hostname"`
	TelegramId string `yaml:"telegram-id"`
}

func NewUser(name, hostname, telegramId string) (*User, error) {
	user := &User{
		Name:     name,
		Hostname: hostname,
	}

	err := user.SetTelegramId(telegramId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u User) DecryptTelegramId() (string, error) {
	if u.TelegramId == "" {
		return "", nil
	}

	return utils.Decrypt(u.TelegramId, env.GetEncryptionKey())
}

func (u *User) SetTelegramId(telegramId string) error {
	encryptedId, err := utils.Encrypt(telegramId, env.GetEncryptionKey())
	if err != nil {
		return err
	}

	u.TelegramId = encryptedId
	return nil
}

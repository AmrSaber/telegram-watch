package models

type User struct {
	Name        string `yaml:"name"`
	TelegramId  string `yaml:"telegram-id"`
	MachineName string `yaml:"machine-name"`
}

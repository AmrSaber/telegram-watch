package common

import "github.com/joho/godotenv"

var envMap map[string]string

const (
	_ENCRYPTION_KEY = "ENCRYPTION_KEY"
	_BOT_TOKEN      = "BOT_TOKEN"
)

func SetEvn(env string) {
	var err error

	envMap, err = godotenv.Unmarshal(env)
	if err != nil {
		panic(err)
	}
}

func GetEncryptionKey() string {
	encryptionKey := envMap[_ENCRYPTION_KEY]

	if encryptionKey == "" {
		panic("ENCRYPTION_KEY was not set in .env at build")
	}

	return encryptionKey
}

func GetBotTokenKey() string {
	botToken := envMap[_BOT_TOKEN]

	if botToken == "" {
		panic("BOT_TOKEN was not set in .env at build")
	}

	return botToken
}

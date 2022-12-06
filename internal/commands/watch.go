package commands

import (
	"io"

	"github.com/AmrSaber/tw/internal/models"
	"github.com/urfave/cli/v2"
)

func WatchCommand(c *cli.Context) error {
	// config := models.LoadConfig()
	// if config.User == nil {
	// 	return fmt.Errorf("no registered user")
	// }

	// telegramId, err := config.User.DecryptTelegramId()
	// if err != nil || telegramId == "" {
	// 	return fmt.Errorf("no registered telegram id")
	// }

	// args := c.Args().Slice()

	// if len(args) == 0 {
	// 	return fmt.Errorf("you must provide a command to run")
	// }

	// command := strings.Join(args, " ")
	// telegramWriter, err := getTelegramWriter(config, command)
	// if err != nil {
	// 	return err
	// }

	// cmd := exec.Command("bash", "-c", strings.Join(args, " "))
	// cmd.Stdout = telegramWriter
	// cmd.Stderr = os.Stderr
	// cmd.Run()

	return nil
}

func getTelegramWriter(config models.Config, command string) (io.WriteCloser, error) {
	// bot, err := tgbotapi.NewBotAPI(env.GetBotTokenKey())
	// if err != nil {
	// 	return nil, err
	// }

	// stringId, _ := config.User.DecryptTelegramId()
	// id, _ := strconv.Atoi(stringId)
	// msgConfig := tgbotapi.NewMessage(int64(id), "Hello, this is a message to track your command")

	// msg, err := bot.Send(msgConfig)
	// if err != nil {
	// 	return nil, err
	// }

	// return models.NewWatcher(bot, msg), nil
	return nil, nil
}

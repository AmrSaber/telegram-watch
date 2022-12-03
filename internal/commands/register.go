package commands

import (
	"fmt"

	"github.com/AmrSaber/tw/internal/common"
	"github.com/urfave/cli/v2"
)

func RegisterCommand(c *cli.Context) error {
	fmt.Println("key:", common.GetEncryptionKey())
	return nil
}

package utils

import (
	"os"
	"path"
)

func GetConfigFilePath() (string, error) {
	configDirectory, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return path.Join(configDirectory, "tw.yaml"), nil
}

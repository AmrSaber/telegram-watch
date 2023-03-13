package models

import (
	"io"
	"os"
	"time"

	"github.com/AmrSaber/tw/internal/env"
	"gopkg.in/yaml.v2"
)

type Config struct {
	User *User `yaml:"user"`

	Runtime RuntimeConfig `yaml:"-"`
}

type RuntimeConfig struct {
	Quiet bool

	Interval time.Duration
	Timeout  time.Duration
}

func LoadConfig() Config {
	configPath, err := env.GetConfigFilePath()
	if err != nil {
		return Config{}
	}

	configFile, err := os.Open(configPath)
	if err != nil {
		return Config{}
	}

	configFileContent, err := io.ReadAll(configFile)
	if err != nil {
		return Config{}
	}

	configFile.Close()

	var config Config
	yaml.Unmarshal(configFileContent, &config)

	return config
}

func (c Config) Save() error {
	configPath, err := env.GetConfigFilePath()
	if err != nil {
		return err
	}

	configFile, err := os.Create(configPath)
	if err != nil {
		return err
	}

	yamlConfig, _ := yaml.Marshal(c)
	if _, err := configFile.Write(yamlConfig); err != nil {
		return err
	}

	return nil
}

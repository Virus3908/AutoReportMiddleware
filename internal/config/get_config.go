package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func GetConfig() (*ConfigStuct, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, fmt.Errorf("не указана переменная CONFIG_PATH до файла конфигурации")
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла конфигурации: %s", err)
	}

	config := &ConfigStuct{}
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, fmt.Errorf("ошибка разбора YAML: %s", err)
	}

	return config, nil
}
package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func GetConfig() (*ConfigStuct, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, fmt.Errorf("the CONFIG_PATH variable for the configuration file is not specified")
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading configuration file: %s", err)
	}

	config := &ConfigStuct{}
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML: %s", err)
	}

	return config, nil
}
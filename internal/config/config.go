package config

import (
	"main/internal/clients"
	"main/internal/database"
	"main/internal/kafka"
	"main/internal/storage"
)

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ConfigStuct struct {
	DB     database.DBConfig `yaml:"pg"`
	Server ServerConfig      `yaml:"server"`
	S3     storage.S3Config  `yaml:"s3"`
	API    clients.APIConfig `yaml:"api"`
	Kafka  kafka.KafkaConfig `yaml:"kafka"`
}

// а валидировать будем как-то? Или нам тут насрут и мы в рантайме упадем?
// обычно это происходит как-то так:

func (c *ConfigStuct) Validate() (error) {
	// check internal fields like this
	if err := c.DB.Validate(); err != nil {
		return err
	}
	// some another validations here
	return nil
}
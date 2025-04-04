package config

import (
	"main/internal/postgres"
	"main/internal/kafka"
	"main/internal/storage"
)

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ConfigStuct struct {
	DB     postgres.DBConfig `yaml:"pg"`
	Server ServerConfig      `yaml:"server"`
	S3     storage.S3Config  `yaml:"s3"`
	Kafka  kafka.KafkaConfig `yaml:"kafka"`
}

package config

import (
	"main/internal/postgres"
	"main/internal/kafka/producer"
	"main/internal/kafka/consumer"
	"main/internal/storage"
)

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	LogLevel string `yaml:"log_level"`
}

type ConfigStuct struct {
	DB     postgres.DBConfig `yaml:"pg"`
	Server ServerConfig      `yaml:"server"`
	S3     storage.S3Config  `yaml:"s3"`
	Producer  producer.KafkaProducerConfig `yaml:"producer"`
	Consumer  consumer.KafkaConsumerConfig `yaml:"consumer"`
}

package config

import (
	"main/internal/kafka/consumer"
	"main/internal/kafka/producer"
	"main/internal/logger"
	"main/internal/postgres"
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
	Producer  producer.KafkaProducerConfig `yaml:"producer"`
	Consumer  consumer.KafkaConsumerConfig `yaml:"consumer"`
	Logger logger.LoggerConfig `yaml:"logger"`
}

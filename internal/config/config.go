package config

import (
	"main/internal/database"
	"main/internal/clients"
	"main/internal/storage"
)

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type ConfigStuct struct {
	DB     database.DBConfig `yaml:"pg"`
	Server ServerConfig          `yaml:"server"`
	S3     storage.S3Config      `yaml:"s3"`
	API    clients.APIConfig     `yaml:"api"`
}

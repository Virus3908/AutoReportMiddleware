package config

import (
	"main/internal/database"
	"main/internal/storage"
)

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type APIConfig struct {
	BaseUrl string `yaml:"baseurl"`
	Timeout int    `yaml:"timeout"`
}

type ConfigStuct struct {
	DB     database.DBConnection `yaml:"pgconnection"`
	Server ServerConfig          `yaml:"server"`
	S3     storage.S3Config      `yaml:"s3"`
	API    APIConfig             `yaml:"api"`
}

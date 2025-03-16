package config

type DBConnection struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type SercerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type S3Config struct {
	Region string `yaml:"region"`
	Bucket string `yaml:"bucket"`
	Endpoint string `yaml:"endpoint"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Minio bool `yaml:"minio"`
}

type ConfigStuct struct {
	DB DBConnection `yaml:"pgconnection"`
	Server   SercerConfig `yaml:"server"`
	S3 S3Config `yaml:"s3"`
}

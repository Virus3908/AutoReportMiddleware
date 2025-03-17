package config

type DBConnection struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type S3Config struct {
	Region    string `yaml:"region"`
	Bucket    string `yaml:"bucket"`
	Endpoint  string `yaml:"endpoint"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Minio     bool   `yaml:"minio"`
}

type APIConfig struct {
	BaseUrl string `yaml:"baseurl"`
	Timeout int    `yaml:"timeout"`
}

type ConfigStuct struct {
	DB     DBConnection `yaml:"pgconnection"`
	Server ServerConfig `yaml:"server"`
	S3     S3Config     `yaml:"s3"`
	API    APIConfig    `yaml:"api"`
}

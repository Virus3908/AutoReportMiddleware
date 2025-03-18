package storage

type S3Config struct {
	Region    string `yaml:"region"`
	Bucket    string `yaml:"bucket"`
	Endpoint  string `yaml:"endpoint"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Minio     bool   `yaml:"minio"`
}
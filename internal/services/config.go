package services

type APIConfig struct {
	BaseURL string `yaml:"baseurl"`
	Timeout int    `yaml:"timeout"`
}
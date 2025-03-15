package config

type DBConnection struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"dbname"`
}

type ConfigStuct struct {
	DBConfig DBConnection `yaml:"pgconnection"`
}

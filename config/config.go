package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Port string `yaml:"port"`

	Rpc struct {
		Network string `yaml:"network"`
		Addr    string `yaml:"addr"`
	} `yaml:"rpc"`

	Redis struct {
		Addr     string `yaml:"addr"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"redis"`

	Mysql struct {
		Dsn string `yaml:"dsn"`
	} `yaml:"mysql"`

	Postgres struct {
		Dsn string `yaml:"dsn"`
	} `yaml:"postgres"`
}

const (
	defaultConfigFile = "./app.yml"

	defaultPort = ":8080"
)

var conf = AppConfig{
	Port: defaultPort,
}

func init() {
	file := os.Getenv("CONFIG_FILE")

	if file == "" {
		file = defaultConfigFile
	}

	f, err := os.ReadFile(file)
	if err != nil {
		return
	}

	if err := yaml.Unmarshal(f, &conf); err != nil {
		return
	}
}

func GetAppConfig() *AppConfig {
	return &conf
}

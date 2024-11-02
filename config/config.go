package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Port string `yaml:"port"`
}

const (
	configFile = "./app.yml"

	defaultPort = ":8080"
)

var conf = AppConfig{
	Port: defaultPort,
}

func init() {
	f, err := os.ReadFile(configFile)
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

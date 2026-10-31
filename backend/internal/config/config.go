package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServiceConfig  `yaml:"service"`
	Database DatabaseConfig `yaml:"database"`
}

type ServiceConfig struct {
	Port        string `yaml:"port"`
	Timeout     int    `yaml:"timeout"`
	Retries     int    `yaml:"retries"`
	IdleTimeout int    `yaml:"idle_timeout"`
}

type DatabaseConfig struct {
	TTL      int    `yaml:"ttl"`
	Timeout  int    `yaml:"timeout"`
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func InitConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return Config{}, err
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	cfg.Database.Password = os.Getenv("DB_PASSWORD")

	return cfg, nil
}

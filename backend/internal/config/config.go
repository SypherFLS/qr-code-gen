package config

import (
	"fmt"
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
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.SSLMode,
	)
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

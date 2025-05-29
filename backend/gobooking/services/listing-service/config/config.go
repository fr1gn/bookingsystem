package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	MongoURI  string `yaml:"mongo_uri"`
	MongoDB   string `yaml:"mongo_db"`
	RedisAddr string `yaml:"redis_addr"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	return &cfg, err
}

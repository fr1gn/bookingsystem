package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	AuthAddress    string `yaml:"auth_address"`
	BookingAddress string `yaml:"booking_address"`
	ListingAddress string `yaml:"listing_address"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	return &cfg, err
}

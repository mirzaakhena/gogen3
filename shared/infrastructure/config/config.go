package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Port     int    `json:"port,omitempty"`
	Host     string `json:"host,omitempty"`
	Database string `json:"database,omitempty"`
}

func ReadConfig() (*Config, error) {

	bytes, err := os.ReadFile("config.json")
	if err != nil {
		return nil, err
	}

	var cfg Config

	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

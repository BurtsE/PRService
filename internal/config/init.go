package config

import (
	"encoding/json"
	"os"
)

func InitConfig() (*Config, error) {
	data, err := os.ReadFile("./configs/config.json")
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

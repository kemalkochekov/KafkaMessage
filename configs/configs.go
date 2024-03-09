package configs

import (
	"encoding/json"

	"os"

	"github.com/go-playground/validator/v10"
)

const configPath = "./configs/configs.json"

type Config struct {
	Server struct {
		Host string `validate:"required"`
	}
	Kafka struct {
		Brokers []string `json:"brokers"`
		Topic   string   `json:"topic"`
	}
}

func LoadConfig() (c *Config, err error) {
	jsonFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(jsonFile).Decode(&c)
	if err != nil {
		return nil, err
	}

	err = validator.New().Struct(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

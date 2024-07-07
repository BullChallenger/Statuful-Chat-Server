package config

import (
	"github.com/naoina/toml"
	"os"
)

type Config struct {
	DB struct {
		Database string
		URL      string
	}

	Kafka struct {
		URL     string
		GroupID string
	}

	Info struct {
		Port string
	}
}

func NewConfig(path string) *Config {
	config := new(Config)

	if file, err := os.Open(path); err != nil {
		panic(err)
	} else if err = toml.NewDecoder(file).Decode(config); err != nil {
		panic(err)
	} else {
		return config
	}
}

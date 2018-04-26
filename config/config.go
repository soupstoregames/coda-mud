package config

import (
	"github.com/codingconcepts/env"
)

type Config struct {
	Address  string `env:"ADDRESS" default:"0.0.0.0"`
	Port     string `env:"PORT" default:"50050"`
	DataPath string `env:"DATA_PATH" default:"/etc/data"`
}

func Load() (config *Config, err error) {
	config = &Config{}

	if err = env.Set(config); err != nil {
		return nil, err
	}

	return
}

package config

import (
	"github.com/codingconcepts/env"
)

type Config struct {
	Address      string `env:"ADDRESS" default:"0.0.0.0"`
	Port         string `env:"PORT" default:"5555"`
	DataPath     string `env:"DATA_PATH" default:"/etc/data"`
	LoginAddress string `env:"LOGIN_ADDRESS" default:"localhost:50051"`

	DatabaseUser     string `env:"DB_USER" default:"Nick"`
	DatabasePassword string `env:"DB_PASS" default:""`
	DatabaseName     string `env:"DB_NAME" default:"coda"`
}

func Load() (config *Config, err error) {
	config = &Config{}

	if err = env.Set(config); err != nil {
		return nil, err
	}

	return
}

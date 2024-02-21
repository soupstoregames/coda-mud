package config

import (
	"github.com/codingconcepts/env"
)

type Config struct {
	// Address to listen on, generally use default
	Address string `env:"ADDRESS" default:"0.0.0.0"`
	// Port to listen on
	Port string `env:"PORT" default:"5555"`
	// DataPath is the absolute path to the data files
	DataPath string `env:"DATA_PATH" default:"/Users/rinse/work/games/coda-data"`
	// StatePath is where the history of the game is stored
	StatePath string `env:"STATE_PATH" default:"state"`
}

func Load() (*Config, error) {
	config := &Config{}
	if err := env.Set(config); err != nil {
		return nil, err
	}

	return config, nil
}

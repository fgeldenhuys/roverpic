package main

import (
	"fmt"

	"github.com/BurntSushi/toml"

	"roverpic/downloader"
	"roverpic/roverapi"
	"roverpic/server"
)

type Config struct {
	RoverAPI   roverapi.Config
	Downloader downloader.Config
	Server     server.Config
}

var configDefaults = Config{
	RoverAPI: roverapi.Defaults,
}

func DecodeConfig(path string) (*Config, error) {
	config := configDefaults
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, fmt.Errorf("Unable to decode config file: %w", err)
	}
	return &config, nil
}

func (conf *Config) Validate() error {
	if err := conf.RoverAPI.Validate(); err != nil {
		return fmt.Errorf("RoverAPI: %w", err)
	}
	return nil
}

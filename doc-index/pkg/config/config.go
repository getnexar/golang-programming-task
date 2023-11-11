package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/creasty/defaults"
)

var cfg *Config

const (
	ConfJsonPathEnvVar = "CONF_JSON_PATH"
)

type Config struct {
	ServerAddress    string `json:"serverAddress" default:":8080"`
	LogLevel         string `json:"logLevel" default:"debug"`
	IndexDataDir     string `json:"indexDataDir" default:"/tmp/data"`
	MaxSearchResults int    `json:"maxSearchResults" default:"10"`
}

func Default() *Config {
	c := &Config{}
	err := defaults.Set(c)
	if err != nil {
		panic(err)
	}
	return c
}

func Load() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}
	path := os.Getenv(ConfJsonPathEnvVar)
	if path == "" {
		return Default(), nil
	}

	var err error
	cfg, err = FromFile(path)
	if err != nil {
		return nil, err
	}

	err = defaults.Set(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func FromFile(configFilePath string) (*Config, error) {
	file, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	c := Default()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

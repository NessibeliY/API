package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const configFile = "config.yaml"

type Config struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	DBPort     int    `yaml:"db_port"`
	DBName     string `yaml:"db_name"`
	DBUser     string `yaml:"db_user"`
	DBPassword string `yaml:"db_password"`
}

func Load() (*Config, error) {
	config := &Config{}
	rawYaml, err := os.ReadFile(configFile)
	if err != nil {
		return nil, errors.Wrap(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYaml, &config)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshaling yaml")
	}
	return config, nil
}

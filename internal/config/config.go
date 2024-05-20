package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const configFile = "configs/config.yaml"

type Config struct {
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	DBPort        int    `yaml:"db_port"`
	DBName        string `yaml:"db_name"`
	DBUser        string `yaml:"db_user"`
	DBPassword    string `yaml:"db_password"`
	DBHost        string `yaml:"db_host"`
	RedisPort     string `yaml:"redis_port"`
	RedisPassword string `yaml:"redis_password"`
	RedisDB       int    `yaml:"redis_db"`
	RedisHost     string `yaml:"redis_host"`
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

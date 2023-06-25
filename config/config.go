package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	Development = "dev"
	Production  = "prod"
)

type Config struct {
	Server   ServerConfig   `toml:"server"`
	Database DatabaseConfig `toml:"database"`
}

type ServerConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type DatabaseConfig struct {
	Name     string `toml:"name"`
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

func LoadConfig() *Config {
	env := os.Getenv("env")

	if env == "" {
		env = Development
	}

	if env != Development && env != Production {
		panic(fmt.Errorf("invalid environment, possible values are: '%s' or '%s'", Development, Production))
	}

	viper.SetConfigName(env)
	viper.SetConfigType("toml")
	viper.AddConfigPath("config/")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error reading config file: %w", err))
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		panic(fmt.Errorf("unable to decode into Config struct: %w", err))
	}

	fmt.Printf("Loaded '%s' config: %+v\n", env, config)
	return config
}

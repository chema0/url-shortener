package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/chema0/url-shortener/pkg/utils"
	"github.com/spf13/viper"
)

const (
	Development = "dev"
	Production  = "prod"
	Test        = "test"
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
	env := utils.Get("env", Development)

	if env != Development && env != Production && env != Test {
		panic(fmt.Errorf("invalid environment, possible values are: '%s', %s or '%s'", Development, Production, Test))
	}

	projectConfigPath, err := getProjectConfigPath()
	if err != nil {
		panic(fmt.Errorf("fatal error reaching config dir: %w", err))
	}

	viper.SetConfigName(env)
	viper.SetConfigType("toml")
	viper.AddConfigPath(projectConfigPath)

	err = viper.ReadInConfig()
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

func getProjectConfigPath() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}

	projectConfigPath := ""
	err := filepath.Walk(filepath.Dir(filename), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && filepath.Base(path) == "config" {
			projectConfigPath = path
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if projectConfigPath == "" {
		return "", fmt.Errorf("failed to find project root path")
	}

	return projectConfigPath, nil
}

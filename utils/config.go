package utils

import (
	"fmt"

	"github.com/caarlos0/env/v7"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	ServerPort int `env:"PORT,required"`

	DbHost     string `env:"DB_MONGO_HOST,required"`
	DbUser     string `env:"DB_MONGO_USER,required"`
	DbName     string `env:"DB_MONGO_NAME,required"`
	DbPort     int    `env:"DB_MONGO_PORT,required"`
	DbPassword string `env:"DB_MONGO_PASSWORD,required"`
}

func LoadConnectionConfig() (*EnvConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf(color.RedString("unable to load .env file: %v"), err)
	}

	config := EnvConfig{}

	err = env.Parse(&config)
	if err != nil {
		return nil, fmt.Errorf(color.RedString("unable to parse environment variables: %v"), err)
	}

	return &config, nil
}

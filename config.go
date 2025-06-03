package main

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AppPort   string `env:"APP_PORT" env-default:"8000"`
	Env       string `env:"ENV" env-default:"local"`
	LogPath   string `env:"LOG_PATH" env-default:""`
	DBUrl     string `env:"DB_URL" env-required:"true"`
	JWTSecret string `env:"JWT_SECRET" env-required:"true"`
}

func MustLoadConfig() Config {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic("Configuration error: " + err.Error())
	}
	return cfg
}

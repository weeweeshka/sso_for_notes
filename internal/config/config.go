package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storagePath" env-default:"postgres://postgres:12345@localhost:5433/notes?sslmode=disable"`
	GRPCServer
}

type GRPCServer struct {
	Port    string `yaml:"port" env-default:"5445"`
	Timeout string `yaml:"timeout" env-default:"5s"`
}

func MustLoad() *Config {
	configPath := "D:/go Projects/notes_auth/sso/config/local.yaml"
	if _, err := os.Stat("D:/go Projects/notes_auth/sso/config/local.yaml"); os.IsNotExist(err) {
		panic("config file not found")
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		panic("failed to read config:")
	}

	return &cfg
}

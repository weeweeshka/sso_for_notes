package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storagePath" env-default:"postgres://postgres:12345@localhost:5430/sso?sslmode=disable"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-default:"1h"`
	GRPCServer
}

type GRPCServer struct {
	Port    int    `yaml:"port" env-default:"5445"`
	Timeout string `yaml:"timeout" env-default:"5s"`
}

func MustLoad() *Config {
	return MustLoadPath("D:/go Projects/notes_auth/sso_for_notes/config/local.yaml")
}

func MustLoadPath(configPath string) *Config {

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file not found: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

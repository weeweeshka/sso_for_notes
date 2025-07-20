package main

import (
	"log/slog"
	"os"
	"sso/internal/config"
)

func main() {
	cfg := config.MustLoad()
	slog.Info("Config loaded")

	slogger := SetupLogger()
	slog.Info("Logger initialized")

	//TODO init app

	//TODO gracefull down

}

func SetupLogger() *slog.Logger {
	var log *slog.Logger

	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	return log
}

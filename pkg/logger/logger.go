package logger

import (
	"io"
	"log/slog"
	"os"
)

type Config struct {
	AppENV  string
	LogPath string
	Debug   bool
}

func SetupLogger(config *Config) {
	var logger *slog.Logger
	err := os.MkdirAll(config.LogPath, 0755)
	if err != nil {
		panic(err)
	}
	file, err := os.OpenFile(config.LogPath+"/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	//defer file.Close()

	var logLevel slog.Level
	if config.AppENV == "production" {
		logLevel = slog.LevelWarn
	} else {
		if config.Debug {
			logLevel = slog.LevelDebug
		} else {
			logLevel = slog.LevelInfo
		}
	}
	logJSONHandler := slog.NewJSONHandler(io.MultiWriter(os.Stdout, file), &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	})
	logger = slog.New(logJSONHandler)
	slog.SetDefault(logger)

	if config.AppENV == "production" {
		slog.Info("log to file")
	} else {
		slog.Info("by default log to stdout, set environment APP_ENV=production to log to file")
	}
}

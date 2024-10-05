package config

import (
	"books-api/pkg/xvalidator"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	AppEnvConfig   *AppConfig
	DatabaseConfig *DatabaseConfig
	AuthConfig     *Auth
}

func (c Config) IsStaging() bool {
	return c.AppEnvConfig.AppEnv != "production"
}

func (c Config) AppName() string {
	return c.AppEnvConfig.AppName
}

func (c Config) IsProd() bool {
	return c.AppEnvConfig.AppEnv == "production"
}

func (c Config) IsDebug() bool {
	return c.AppEnvConfig.AppDebug
}

func InitAppConfig(validate *xvalidator.Validator) *Config {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		slog.Error(fmt.Sprintf("Failed to read config file: %s", err))
		slog.Info(fmt.Sprintf("Reading from default environment...: %s", err))
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	viper.WatchConfig()

	// Iterate through all the settings
	for key, value := range viper.AllSettings() {
		envKey := strings.ToUpper(key)
		err := os.Setenv(strings.ToUpper(key), fmt.Sprintf("%v", value))
		if err != nil {
			slog.Error(fmt.Sprintf("Error setting environment variable %s: %v", envKey, err))
			os.Exit(1)
		}
	}
	c := Config{
		AppEnvConfig:   AppConfigInit(),
		DatabaseConfig: DatabaseConfigConfig(),
		AuthConfig:     AuthConfig(),
	}
	errs := validate.Struct(c)
	if errs != nil {
		for k, e := range errs {
			slog.Error(fmt.Sprintf("Failed to load env: %s, msg: %s", k, strings.ToLower(e)))
		}
		os.Exit(1)
	}
	slog.Info("Config loaded")
	return &c
}

package config

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	AppEnv       string
	AppDebug     bool
	AppVersion   string   `validate:"required,startswith=v,alphanum" name:"APP_VERSION"`
	AppName      string   `validate:"required" name:"APP_NAME"`
	HttpPort     string   `validate:"required,number" name:"HTTP_PORT"`
	AllowOrigins []string `name:"HTTP_ALLOW_ORIGINS"`
	AllowMethods []string `name:"HTTP_ALLOW_METHODS"`
	AllowHeaders []string `name:"HTTP_ALLOW_HEADERS"`
	UseReplica   bool     `validate:"boolean" name:"USE_REPLICA"`
	LogFilePath  string   `validate:"required" name:"LOG_PATH"`
}

func AppConfigInit() *AppConfig {
	return &AppConfig{
		AppEnv:       viper.GetString("APP_ENV"),
		AppDebug:     viper.GetBool("APP_DEBUG"),
		AppVersion:   viper.GetString("APP_VERSION"),
		AppName:      viper.GetString("APP_NAME"),
		HttpPort:     viper.GetString("HTTP_PORT"),
		UseReplica:   viper.GetBool("USE_REPLICA"),
		LogFilePath:  viper.GetString("LOG_PATH"),
		AllowOrigins: viper.GetStringSlice("ALLOW_ORIGINS"),
		AllowMethods: viper.GetStringSlice("ALLOW_METHODS"),
		AllowHeaders: viper.GetStringSlice("ALLOW_HEADERS"),
	}
}

package config

import (
	"github.com/spf13/viper"
)

type Auth struct {
	JwtSecretAccessToken string `validate:"required" name:"JWT_SECRET_ACCESS_TOKEN"`
}

func AuthConfig() *Auth {
	return &Auth{
		JwtSecretAccessToken: viper.GetString("JWT_SECRET_ACCESS_TOKEN"),
	}
}

package config

import (
	"github.com/spf13/viper"
)

type RedisConfig struct {
	Redisdatabase int    `validate:"number,min=0" name:"REDIS_DB"`
	Redishost     string `validate:"required" name:"REDIS_HOST"`
	Redisport     int    `validate:"required,number" name:"REDIS_PORT"`
	Redispassword string `name:"REDIS_PASSWORD"`
}

func RedisConfigInit() *RedisConfig {
	return &RedisConfig{
		Redisdatabase: viper.GetInt("REDIS_DB"),
		Redishost:     viper.GetString("REDIS_HOST"),
		Redisport:     viper.GetInt("REDIS_PORT"),
		Redispassword: viper.GetString("REDIS_PASSWORD"),
	}
}

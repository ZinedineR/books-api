package config

import (
	"github.com/spf13/viper"
)

type RabbiterConfig struct {
	RabbitMQDial string `validate:"required" name:"RABBITMQ_DIAL"`
}

func RabbiterConfigInit() *RabbiterConfig {
	return &RabbiterConfig{
		RabbitMQDial: viper.GetString("RABBITMQ_DIAL"),
	}
}

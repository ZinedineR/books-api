package config

import (
	"github.com/spf13/viper"
)

type PubSubConfig struct {
	PubSubService string `validate:"required,eq=redis|eq=kafka" name:"PUBSUB_SERVICE"`
}

func PubSubConfigInit() *PubSubConfig {
	return &PubSubConfig{
		PubSubService: viper.GetString("PUBSUB_SERVICE"),
	}
}

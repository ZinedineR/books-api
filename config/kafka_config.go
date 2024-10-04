package config

import (
	"github.com/spf13/viper"
)

type KafkaConfig struct {
	KafkaSecurityProtocol  string   `validate:"required,eq=SASL_SSL|eq=SASL_PLAIN|eq=PLAIN" name:"KAFKA_SECURITY_PROTOCOL"`
	KafkaUsername          string   `validate:"required" name:"KAFKA_USERNAME"`
	KafkaPassword          string   `validate:"required" name:"KAFKA_PASSWORD"`
	KafkaBroker            []string `validate:"required" name:"KAFKA_BROKERS"`
	KafkaGroupId           string   `validate:"required" name:"KAFKA_GROUP_ID"`
	KafkaTopicNotification string   `validate:"required" name:"KAFKA_TOPIC_NOTIFICATION"`
	KafkaTopicEmail        string   `validate:"required" name:"KAFKA_TOPIC_EMAIL"`
	KafkaTopicSMS          string   `validate:"required" name:"KAFKA_TOPIC_SMS"`
	KafkaTopicGotify       string   `validate:"required" name:"KAFKA_TOPIC_GOTIFY"`
}

func KafkaConfigInit() *KafkaConfig {
	return &KafkaConfig{
		KafkaSecurityProtocol:  viper.GetString("KAFKA_SECURITY_PROTOCOL"),
		KafkaUsername:          viper.GetString("KAFKA_USERNAME"),
		KafkaPassword:          viper.GetString("KAFKA_PASSWORD"),
		KafkaBroker:            viper.GetStringSlice("KAFKA_BROKERS"),
		KafkaGroupId:           viper.GetString("KAFKA_GROUP_ID"),
		KafkaTopicNotification: viper.GetString("KAFKA_TOPIC_NOTIFICATION"),
		KafkaTopicEmail:        viper.GetString("KAFKA_TOPIC_EMAIL"),
		KafkaTopicSMS:          viper.GetString("KAFKA_TOPIC_SMS"),
		KafkaTopicGotify:       viper.GetString("KAFKA_TOPIC_GOTIFY"),
	}
}

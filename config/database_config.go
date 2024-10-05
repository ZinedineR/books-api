package config

import (
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Dbservice  string `validate:"required,eq=postgres|eq=mysql|eq=sqlserver|eq=sqlite" name:"DB_CONNECTION"`
	Dbhost     string `name:"DB_HOST"`
	Dbport     int    `name:"DB_PORT"`
	Dbname     string `name:"DB_DATABASE"`
	Dbuser     string `name:"DB_USERNAME"`
	Dbpassword string `name:"DB_PASSWORD"`
	DbPrefix   string `validate:"required" name:"DB_PREFIX"`
}

func DatabaseConfigConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Dbservice:  viper.GetString("DB_CONNECTION"),
		Dbhost:     viper.GetString("DB_HOST"),
		Dbport:     viper.GetInt("DB_PORT"),
		Dbname:     viper.GetString("DB_DATABASE"),
		Dbuser:     viper.GetString("DB_USERNAME"),
		Dbpassword: viper.GetString("DB_PASSWORD"),
		DbPrefix:   viper.GetString("DB_PREFIX"),
	}
}

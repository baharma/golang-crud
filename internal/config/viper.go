package config

import (
	"log"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	v := viper.New()
	v.SetConfigFile(".env")

	if err := v.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %v", err)
	}

	return v
}

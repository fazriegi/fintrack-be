package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var viperInstance *viper.Viper

func NewViper() *viper.Viper {
	config := viper.New()

	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("../../")
	config.AddConfigPath("./")
	err := config.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	viperInstance = config

	return config
}

func GetString(key string) string {
	return viperInstance.GetString(key)
}

func GetInt(key string) int {
	return viperInstance.GetInt(key)
}

func GetUint(key string) uint {
	return viperInstance.GetUint(key)
}

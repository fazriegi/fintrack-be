package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var viperInstance *viper.Viper

func NewViper() {
	config := viper.New()

	config.SetConfigFile(".env")
	config.SetConfigType("env")

	config.SetEnvPrefix("APP")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.AutomaticEnv()

	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	viperInstance = config
	fmt.Println("JWT_SECRET:", viperInstance.GetString("JWT_SECRET"))
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

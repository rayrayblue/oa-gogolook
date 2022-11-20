package domain

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string, configName string) (config AppConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(configName)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

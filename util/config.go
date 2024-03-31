package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBSource    string `mapstructure:"DB_SOURCE""`
	HostAddress string `mapstructure:"HOST_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

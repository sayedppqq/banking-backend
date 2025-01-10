package util

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	DBSource             string        `mapstructure:"DB_SOURCE"`
	MigrationsDirectory  string        `mapstructure:"MIGRATIONS_DIRECTORY"`
	HTTPHostAddress      string        `mapstructure:"HTTP_HOST_ADDRESS"`
	GrpcHostAddress      string        `mapstructure:"GRPC_HOST_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
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

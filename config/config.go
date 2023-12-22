package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBURL                string        `mapstructure:"PG_URL"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func NewConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)

	viper.SetConfigName("env")
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)

	viper.SetConfigType("env")

	return
}

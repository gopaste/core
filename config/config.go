package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBURL                 string `mapstructure:"PG_URL"`
	AccessTokenSecret     string `mapstructure:"ACCESS_TOKEN_SECRET"`
	AccessTokenExpiryHour int    `mapstructure:"ACCESS_TOKEN_EXPIRY_HOUR"`
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

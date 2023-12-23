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
	AWSSecretKey         string        `mapstructure:"AWS_SECRET_KEY"`
	AWSAccessKey         string        `mapstructure:"AWS_ACCESS_KEY"`
	AWSRegion            string        `mapstructure:"AWS_REGION"`
	AWSSenderEmail       string        `mapstructure:"AWS_SENDER_EMAIL"`
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

package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBUrl             string `mapstructure:"DB_URL"`
	RedisUrl          string `mapstructure:"REDIS_URL"`
	JWTExpireHours    int    `mapstructure:"JWT_EXPIRE_HOURS"`
	JWTIssuer         string `mapstructure:"JWT_ISSUER"`
	JWTTokenSignature string `mapstructure:"JWT_TOKEN_SIGNATURE"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	// TODO: Handle production/beta/develop.
	viper.SetConfigName(".development")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

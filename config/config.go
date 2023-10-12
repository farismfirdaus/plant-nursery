package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port     int            `mapstructure:"PORT"`
	Database DatabaseConfig `mapstructure:",squash"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Database string `mapstructure:"DB_NAME"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
}

func InitConfig() *Config {
	viper.SetConfigName("env")
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	config := &Config{}

	err = viper.Unmarshal(config)
	if err != nil {
		panic(err)
	}

	return config
}

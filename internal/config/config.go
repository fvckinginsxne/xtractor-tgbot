package config

import (
	"bot/pkg/tech/e"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Hostname    string `mapstructure:"hostname"`
	StoragePath string `mapstructure:"storage_path"`
}

func Init() (*Config, error) {
	viper.AddConfigPath("/Users/madw3y/petprojects/spotik-bot/configs")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, e.Wrap("can't read config", err)
	}

	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, e.Wrap("can't unmarshal config", err)
	}

	if err := parseEnv(); err != nil {
		return nil, e.Wrap("can't parse env", err)
	}

	return &cfg, nil
}

func parseEnv() error {
	if err := godotenv.Load("/Users/madw3y/petprojects/spotik-bot/.env"); err != nil {
		return e.Wrap("can't load env file", err)
	}

	return nil
}

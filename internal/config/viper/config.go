package viper

import (
	"fmt"

	"github.com/AlisskaPie/project-xm/internal/config"

	"github.com/spf13/viper"
)

func GetConfig() (config.Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return config.Config{}, fmt.Errorf("failed to read in config: %w", err)
	}

	conf := config.Config{}
	if err := viper.Unmarshal(&conf); err != nil {
		return config.Config{}, fmt.Errorf("failed to unmarshal config into struct: %w", err)
	}

	return conf, nil
}

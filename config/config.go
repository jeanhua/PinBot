package config

import (
	"github.com/spf13/viper"
	"log"
)

// 配置
var config *viper.Viper

func GetConfig() *viper.Viper {
	return config
}

func LoadConfig() {
	cfg := viper.New()
	cfg.SetConfigName("config")
	cfg.AddConfigPath(".")
	cfg.SetConfigType("yaml")
	if err := cfg.ReadInConfig(); err != nil {
		log.Fatalf("Failed to load config file: %v", err)
	}
	config = cfg
}

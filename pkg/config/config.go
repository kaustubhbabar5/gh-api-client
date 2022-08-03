package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// holds the app configuration
type Config struct {
	LogLevel        string `mapstructure:"LOG_LEVEL"`
	LogTimeFormat   string `mapstructure:"LOG_TIME_FORMAT"`
	Host            string `mapstructure:"HOST"`
	Port            string `mapstructure:"PORT"`
	RedisUrl        string `mapstructure:"REDIS_URL"`
	RedisPassword   string `mapstructure:"REDIS_PASSWORD"`
	ReadTimeout     int    `mapstructure:"READ_TIMEOUT"`
	WriteTimeout    int    `mapstructure:"WRITE_TIMEOUT"`
	GithubAuthToken string `mapstructure:"GITHUB_AUTH_TOKEN"`
}

// discovers and loads the configuration file
func Load(path string, logger *zap.Logger) (*Config, error) {
	setDefaults()

	viper.SetConfigName(".env")
	viper.SetConfigType("dotenv")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if !ok {
			return nil, fmt.Errorf("viper.ReadInConfig :%w", err)
		}
		logger.Warn("config not found", zap.Error(err))
	}

	var config Config

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("viper.Unmarshal :%w", err)
	}

	return &config, nil
}

// loads default config
func setDefaults() {
	viper.SetDefault("LOG_LEVEL", zap.InfoLevel)
	viper.SetDefault("LOG_TIME_FORMAT", time.RFC3339Nano)

	viper.SetDefault("HOST", "0.0.0.0")
	viper.SetDefault("PORT", "8080")
}

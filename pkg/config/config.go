package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// holds the app configuration.
type Config struct {
	LogLevel           string `mapstructure:"LOG_LEVEL"`
	LogTimeFormat      string `mapstructure:"LOG_TIME_FORMAT"`
	Host               string `mapstructure:"HOST"`
	Port               string `mapstructure:"PORT"`
	RedisURL           string `mapstructure:"REDIS_URL"`
	RedisPasswordKey   string `mapstructure:"REDIS_PASSWORD_ENV_KEY"`
	ReadTimeout        int    `mapstructure:"READ_TIMEOUT"`
	WriteTimeout       int    `mapstructure:"WRITE_TIMEOUT"`
	ReadHeaderTimeout  int    `mapstructure:"READ_HEADER_TIMEOUT"`
	GithubAuthTokenKey string `mapstructure:"GITHUB_AUTH_TOKEN_ENV_KEY"`
}

// discovers and loads the configuration file in given path.
func Load(path string, logger *zap.Logger) (*Config, error) {
	setDefaults()

	err := godotenv.Load()
	if err != nil {
		logger.Debug("dotenv not found", zap.Error(err))
	}

	viper.SetConfigName(".env")
	viper.SetConfigType("dotenv")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		var notFoundError *viper.ConfigFileNotFoundError
		if errors.As(err, &notFoundError) {
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

// loads default config.
func setDefaults() {
	viper.SetDefault("LOG_LEVEL", zap.InfoLevel)
	viper.SetDefault("LOG_TIME_FORMAT", time.RFC3339Nano)

	viper.SetDefault("HOST", "0.0.0.0")
	viper.SetDefault("PORT", "8080")

	viper.SetDefault("REDIS_URL", "0.0.0.0:6379")
	viper.SetDefault("REDIS_PASSWORD_ENV_KEY", "REDIS_PASSWORD")

	viper.SetDefault("WRITE_TIMEOUT", "30")
	viper.SetDefault("READ_TIMEOUT", "30")
	viper.SetDefault("READ_HEADER_TIMEOUT", "2")

	viper.SetDefault("GITHUB_AUTH_TOKEN_ENV_KEY", "GITHUB_AUTH_TOKEN")
}

package config

import "os"

type Config struct {
	HTTPAddress string
	LogLevel    string
}

func Load() *Config {
	return &Config{
		HTTPAddress: getEnv("HTTP_ADDRESS", ":8080"),
		LogLevel:    getEnv("LOG_LEVEL", "INFO"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

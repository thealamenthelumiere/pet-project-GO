package configs

import (
	"os"
	"strconv"
)

type Config struct {
	Port   int
	Secret string
}

func LoadConfig() *Config {
	port := getEnvAsInt("PORT", 8080)
	secret := getEnv("SECRET_KEY", "default-secret-key-change-in-production")

	return &Config{
		Port:   port,
		Secret: secret,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

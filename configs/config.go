package configs

import (
    "time"

    "github.com/joho/godotenv"
    "github.com/spf13/viper"
)

type Config struct {
    Port            int           `mapstructure:"PORT"`
    Secret          string        `mapstructure:"SECRET_KEY"`
    AccessTokenTTL  time.Duration `mapstructure:"ACCESS_TOKEN_TTL"`
    RefreshTokenTTL time.Duration `mapstructure:"REFRESH_TOKEN_TTL"`
}

func LoadConfig() (*Config, error) {
    _ = godotenv.Load()

    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("./configs")

    viper.SetDefault("PORT", 8080)
    viper.SetDefault("SECRET_KEY", "default-secret-key-change-in-production")
    viper.SetDefault("ACCESS_TOKEN_TTL", "15m")
    viper.SetDefault("REFRESH_TOKEN_TTL", "24h")

    viper.AutomaticEnv()
    _ = viper.ReadInConfig()

    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}
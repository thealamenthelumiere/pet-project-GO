package configs

import (
	
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config содержит все настройки приложения.
// Теги mapstructure нужны для корректного маппинга из Viper.
type Config struct {
	Port   int    `mapstructure:"PORT"`
	Secret string `mapstructure:"SECRET_KEY"`
}

// LoadConfig загружает конфигурацию из .env, файлов и переменных окружения.
// Возвращает указатель на Config и ошибку (если не удалось распарсить).
func LoadConfig() (*Config, error) {
	// 1. Загружаем .env файл (если есть) в окружение.
	// Ошибка игнорируется – файл может отсутствовать.
	_ = godotenv.Load()

	// 2. Настраиваем Viper.
	viper.SetConfigName("config")      // имя файла без расширения
	viper.SetConfigType("yaml")        // поддерживаются yaml, json, toml, env и др.
	viper.AddConfigPath(".")           // ищем в текущей директории
	viper.AddConfigPath("./configs")   // или в подпапке configs

	// 3. Значения по умолчанию.
	viper.SetDefault("PORT", 8080)
	viper.SetDefault("SECRET_KEY", "default-secret-key-change-in-production")

	// 4. Автоматически читать переменные окружения.
	viper.AutomaticEnv()

	// 5. Пытаемся прочитать файл конфигурации (необязательно).
	// Если файла нет – работаем с дефолтами и переменными окружения.
	_ = viper.ReadInConfig()

	// 6. Десериализуем всё в структуру Config.
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

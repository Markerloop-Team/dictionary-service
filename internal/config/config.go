package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config конфигурация приложения
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig    `json:"redis"`
	Log      LogConfig      `json:"log"`
	Cache    CacheConfig    `json:"cache"`
}

// CacheConfig конфигурация кэша
type CacheConfig struct {
	TTL time.Duration
}

// ServerConfig конфигурация сервера
type ServerConfig struct {
	Host         string        `json:"host"`
	Port         string        `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
}

// DatabaseConfig конфигурация базы данных
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int           // Максимальное количество открытых соединений
	MaxIdleConns    int           // Максимальное количество неактивных соединений
	ConnMaxLifetime time.Duration // Максимальное время жизни соединения
}

// RedisConfig конфигурация Redis
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// LogConfig конфигурация логирования
type LogConfig struct {
	Level string `json:"level"`
}

// Load загружает конфигурацию из переменных окружения
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8084"),
			Host:         getEnv("HOST", "0.0.0.0"),
			ReadTimeout:  getEnvDuration("READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getEnvDuration("WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getEnvDuration("IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			Host:            getEnvRequired("DB_HOST"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnvRequired("DB_USER"),
			Password:        getEnvRequired("DB_PASSWORD"),
			Name:            getEnvRequired("DB_NAME"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			Addr:     getEnvRequired("REDIS_ADDR"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 2), // DB 2 для dictionary-service
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
		Cache: CacheConfig{
			TTL: getEnvDuration("CACHE_TTL", 1*time.Hour),
		},
	}
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvDuration получает значение переменной окружения как duration
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// getEnvInt получает целочисленное значение переменной окружения
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvRequired получает обязательную переменную окружения или паникует
func getEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("обязательная переменная окружения %s не установлена", key))
	}
	return value
}

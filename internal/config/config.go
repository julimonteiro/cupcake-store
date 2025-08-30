package config

import "os"

type Config struct {
	Port, DBDialect, DBDSN, LogLevel string
}

func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		DBDialect: getEnv("DB_DIALECT", "sqlite"),
		DBDSN:     getEnv("DB_DSN", "cupcake_store.db"),
		LogLevel:  getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

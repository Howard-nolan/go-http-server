package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port    int
	DbPath  string
	BaseURL string
}

func Load() Config {
	return Config{
		Port:    envInt("PORT", 8080),
		DbPath:  envStr("DB_PATH", "/tmp/app.db"),
		BaseURL: envStr("BASE_URL", "http://localhost:8080"),
	}
}

func envInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return defaultVal
}

func envStr(key string, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

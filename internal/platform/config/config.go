package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port int
}

func Load() Config {
	return Config{
		Port: envInt("PORT", 8080),
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

package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port       string
	DBPath     string
	JWTSecret  string
	JWTTTL     int // 秒
	RefreshTTL int // 秒
}

func Load() *Config {
	cfg := &Config{
		Port:       getEnv("PORT", "8080"),
		DBPath:     getEnv("DB_PATH", "data/app.db"),
		JWTSecret:  getEnv("JWT_SECRET", "go-vue-admin-jwt-secret"),
		JWTTTL:     7200,
		RefreshTTL: 7 * 24 * 3600,
	}
	return cfg
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func (c *Config) Address() string {
	return fmt.Sprintf(":%s", c.Port)
}

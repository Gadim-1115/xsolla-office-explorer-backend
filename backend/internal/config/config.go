package config

import (
	"os"
	"time"
)

type Config struct {
	HTTPPort        string
	DatabaseURL     string
	AllowedOrigin   string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration

	// LMSJWTSecret is the LMS backend's JWT_ACCESS_SECRET. Empty: LMS sign-in is off, everyone is a guest.
	LMSJWTSecret string
}

func Load() *Config {
	return &Config{
		HTTPPort:        getEnv("HTTP_PORT", "8080"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://postgres:12345@localhost:5432/xsolla_explorer?sslmode=disable"),
		AllowedOrigin:   getEnv("ALLOWED_ORIGIN", "http://localhost:5173"),
		LMSJWTSecret:    getEnv("LMS_JWT_SECRET", ""),
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		ShutdownTimeout: 10 * time.Second,
	}
}

func getEnv(key, fallback string) string {
	if val, exists := os.LookupEnv(key); exists && val != "" {
		return val
	}
	return fallback
}

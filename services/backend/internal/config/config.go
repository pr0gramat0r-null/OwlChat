package config

import "os"

type Config struct {
	HTTPAddr string
	JWTSecret string
}

func FromEnv() Config {
	return Config{
		HTTPAddr: getenv("HTTP_ADDR", ":8080"),
		JWTSecret: getenv("JWT_SECRET", "dev-secret-change-me"),
	}
}

func getenv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

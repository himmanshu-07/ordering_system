package config

import (
	"os"
)

type Config struct {
	Port           string
	DatabaseURL    string
	BaseURL        string
	JWTSecret      string
	UPIID          string
	WhatsAppNumber string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://orangemango:binarySearch_3@localhost:5432/qr_ordering?sslmode=disable"),
		BaseURL:        getEnv("BASE_URL", "http://localhost:8080"),
		JWTSecret:      getEnv("JWT_SECRET", ""),
		UPIID:          getEnv("UPI_ID", "9636998667-2@ybl"),
		WhatsAppNumber: getEnv("WHATSAPP_NUMBER", "917579428375"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

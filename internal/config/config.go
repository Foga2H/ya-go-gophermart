package config

import (
	"flag"
	"os"
)

var (
	flagRunBaseURL          = flag.String("a", "localhost:8080", "address and port to run server")
	jwtSecret               = flag.String("s", "dev-secret", "secret key for user cookie signing")
	accrualSystemAddressURL = flag.String("r", "localhost:8080", "address and port for accrual system")
)

type Config struct {
	BaseURL                 string
	JWTSecret               string
	AccrualSystemAddressURL string
}

func NewConfig() *Config {
	if envRunBaseURL := os.Getenv("RUN_ADDRESS"); envRunBaseURL != "" {
		*flagRunBaseURL = envRunBaseURL
	}

	if envJWTSecret := os.Getenv("JWT_SECRET"); envJWTSecret != "" {
		*jwtSecret = envJWTSecret
	}

	if envAccrualSystemAddressURL := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddressURL != "" {
		*accrualSystemAddressURL = envAccrualSystemAddressURL
	}

	return &Config{
		BaseURL:                 *flagRunBaseURL,
		JWTSecret:               *jwtSecret,
		AccrualSystemAddressURL: *accrualSystemAddressURL,
	}
}

func NewConfigFrom(baseURL string) *Config {
	return &Config{
		BaseURL:                 baseURL,
		JWTSecret:               "test-secret",
		AccrualSystemAddressURL: *accrualSystemAddressURL,
	}
}

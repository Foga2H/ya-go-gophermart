package db

import (
	"flag"
	"os"
)

var (
	flagDatabaseDsn = flag.String("d", "", "database connection string")
)

type Config struct {
	DatabaseDSN string
}

func NewConfig() *Config {
	if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
		*flagDatabaseDsn = envDatabaseURI
	}

	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" {
		*flagDatabaseDsn = envDatabaseDSN
	}

	return &Config{
		DatabaseDSN: *flagDatabaseDsn,
	}
}

func NewConfigFrom(databaseDSN string) *Config {
	return &Config{
		DatabaseDSN: databaseDSN,
	}
}

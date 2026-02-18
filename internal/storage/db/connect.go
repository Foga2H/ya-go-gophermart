package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"

	dbConfig "github.com/Foga2H/ya-go-gophermart/internal/config/db"
)

var ErrPingDatabase = errors.New("error pinging database")

type DatabaseConnection struct {
	*dbConfig.Config
}

func NewDatabaseConnection(dbConfig *dbConfig.Config) *DatabaseConnection {
	return &DatabaseConnection{dbConfig}
}

func (conn *DatabaseConnection) Ping() (*sql.DB, error) {
	dbConnection, err := sql.Open("pgx", conn.Config.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	if err := dbConnection.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPingDatabase, err)
	}

	return dbConnection, nil
}

package postgres

import (
	"context"
	"fmt"

	"github.com/ekifel/playlist/internal/config"
	"github.com/jackc/pgx/v5"
)

func NewPostgresClient(cfg *config.PostgresConfig) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to database: %v\n", err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error occurred while pinging postgres: %v", err.Error())
	}

	return conn, nil
}

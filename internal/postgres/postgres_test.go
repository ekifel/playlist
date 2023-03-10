package postgres_test

import (
	"os"
	"testing"

	"github.com/ekifel/playlist/internal/config"
	"github.com/ekifel/playlist/internal/postgres"
)

func TestNewPostgresClient(t *testing.T) {
	cfg := &config.PostgresConfig{DatabaseURL: os.Getenv("DATABASE_URL")}
	conn, err := postgres.NewPostgresClient(cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if conn == nil {
		t.Error("expected a non-nil connection object")
	}
}

func TestNewPostgresClient_InvalidConfig(t *testing.T) {
	cfg := &config.PostgresConfig{DatabaseURL: "invalid-url"}
	conn, err := postgres.NewPostgresClient(cfg)
	if err == nil {
		t.Error("expected an error")
	}
	if conn != nil {
		t.Error("expected a nil connection object")
	}
}

package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/lgu-elo/auth/internal/config"
)

func New(cfg *config.Cfg) *pgx.Conn {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)

	conn, _ := pgx.Connect(context.Background(), dsn)

	return conn
}

//go:build migrate
// +build migrate

package app

import (
	"fmt"
	"os"

	"github.com/lgu-elo/auth/internal/config"
	"github.com/lgu-elo/auth/pkg/migrate"
)

func init() {
	cfg := config.New()
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)

	migrate.Do(migrate.Down, "./migrations", dsn)
	migrate.Do(migrate.Up, "./migrations", dsn)
	os.Exit(0)
}

package app

import (
	"fmt"
	"log"
	"marcyHomeService/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(cfg *config.Config) {
	m, err := migrate.New(
		"file://database/migration",
		fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.PostgresSQL.Username,
			cfg.PostgresSQL.Password,
			cfg.PostgresSQL.Host,
			cfg.PostgresSQL.Port,
			cfg.PostgresSQL.Database,
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}

package postgresql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Client interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginFunc(ctx context.Context, f func(pgx.Tx) error) error
	BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(pgx.Tx) error) error
	Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, arguments ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

type pgConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func NewPgConfig(username, password, host, port, database string) *pgConfig {
	return &pgConfig{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		Database: database,
	}
}

func NewClient(ctx context.Context, maxAttempts int, maxDelay time.Duration, cfg *pgConfig) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	timeout := 5 * time.Second
	err = doWithAttempts(func() error {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		pgxConfig, err := pgxpool.ParseConfig(dsn)

		if err != nil {
			log.Fatalf("unable to parse config: %v\n", err)
		}

		pool, err = pgxpool.ConnectConfig(ctx, pgxConfig)

		if err != nil {
			log.Print("failed to connect to postgres... going to do next attempt", err)

			return err
		}

		return nil
	}, maxAttempts, maxDelay)

	if err != nil {
		log.Fatalf("error do with tries postgresql: %v", err)
	}

	return pool, nil
}

func doWithAttempts(fn func() error, maxAttempts int, delay time.Duration) error {
	var err error

	for maxAttempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			maxAttempts--

			continue
		}

		return nil
	}

	return err
}

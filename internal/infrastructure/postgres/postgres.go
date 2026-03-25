package postgres

import (
	"context"
	"fmt"

	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // registers postgres driver
    _ "github.com/golang-migrate/migrate/v4/source/file"       // registers file source
)


type Postgres struct {
	Pool *pgxpool.Pool
}

func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error){
	// Parse dsn string into config struct that you can customize
	cfg, err := pgxpool.ParseConfig(dsn)
	
	if err != nil {
		return nil, fmt.Errorf("postgres NewPool error: %w", err)
	}
    

	cfg.MaxConns = 30 
	cfg.MinConns = 2 

	// Any connection that is older than an hour will be closed and reestablished
	// Even though the connection is busy, it will still be closed
	cfg.MaxConnLifetime = time.Hour

	// NewWith means create a new connection pool 
	// With a config
	pool, err := pgxpool.NewWithConfig(ctx, cfg)

	if err != nil {
		return nil, fmt.Errorf("postgres NewPool error: %w", err) 
	}

	return pool, nil 
}

func RunMigrations(dsn string) error {
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return fmt.Errorf("infrastructure postgres run migrations function error : %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("infrastructure postgres run migrations function error : %w", err)
	}

	return nil 
}
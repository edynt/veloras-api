package initialize

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/edynnt/veloras-api/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfpg := cfg.PostgreSQL

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfpg.Username, cfpg.Password, cfpg.Host, cfpg.Port, cfpg.Database, cfpg.SslMode)

	fmt.Println("strs:::" + dsn)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgxpool config: %w", err)
	}

	// Customize pool if needed
	poolConfig.MaxConns = 100
	poolConfig.MinConns = 10

	db, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to pgxpool: %w", err)
	}

	// Kiểm tra kết nối
	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping pgxpool: %w", err)
	}

	log.Println("Database connection established successfully with pgxpool.")

	DB = db
	return db, nil
}

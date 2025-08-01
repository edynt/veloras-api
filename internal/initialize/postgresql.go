package initialize

import (
	"context"
	"fmt"
	"time"

	"github.com/edynnt/veloras-api/pkg/config"
	"github.com/edynnt/veloras-api/pkg/global"
	"github.com/edynnt/veloras-api/pkg/response/msg"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var DB *pgxpool.Pool

func InitDB(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfpg := cfg.PostgreSQL

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfpg.Username, cfpg.Password, cfpg.Host, cfpg.Port, cfpg.Database, cfpg.SslMode)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.CannotConnectPgxpool, err)
	}

	// Customize pool if needed
	poolConfig.MaxConns = 100
	poolConfig.MinConns = 10

	db, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.CannotConnectPgxpool, err)
	}

	// Kiểm tra kết nối
	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", msg.CannotPingPgxpool, err)
	}

	global.Logger.Info("Initialize Postgresql successfully!!", zap.String("ok", "success"))

	DB = db
	return db, nil
}

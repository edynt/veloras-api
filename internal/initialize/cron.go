package initialize

import (
	"github.com/edynt/chogiare/veloras-api/internal/cron"
	"github.com/edynt/chogiare/veloras-api/internal/shared/gen"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitCronScheduler(db *pgxpool.Pool) *cron.Scheduler {
	queries := gen.New(db)
	return cron.NewScheduler(queries)
}


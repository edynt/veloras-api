package initialize

import (
	"github.com/edynnt/veloras-api/internal/cron"
	"github.com/edynnt/veloras-api/internal/shared/gen"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitCronScheduler(db *pgxpool.Pool) *cron.Scheduler {
	queries := gen.New(db)
	return cron.NewScheduler(queries)
}


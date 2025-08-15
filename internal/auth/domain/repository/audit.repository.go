package repository

import (
	"context"

	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
)

type AuditRepository interface {
	CreateAuditLog(ctx context.Context, log *entity.AuditLog) (string, error)
	GetAuditLogsByUser(ctx context.Context, userId string, limit int) ([]*entity.AuditLog, error)
	GetAuditLogsByResource(ctx context.Context, resourceType string, resourceId string) ([]*entity.AuditLog, error)
	GetAuditLogsByAction(ctx context.Context, action string, limit int) ([]*entity.AuditLog, error)
	GetAuditLogsByDateRange(ctx context.Context, startDate, endDate int64) ([]*entity.AuditLog, error)
	DeleteOldAuditLogs(ctx context.Context, cutoffDate int64) error
}

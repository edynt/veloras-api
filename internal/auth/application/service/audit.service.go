package service

import (
	"context"

	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
)

type AuditService interface {
	CreateAuditLog(ctx context.Context, req appDto.CreateAuditLogRequest) (string, error)
	GetAuditLogsByUser(ctx context.Context, userId string, limit int) ([]*appDto.AuditLogResponse, error)
	GetAuditLogsByResource(ctx context.Context, resourceType, resourceId string) ([]*appDto.AuditLogResponse, error)
	GetAuditLogsByAction(ctx context.Context, action string, limit int) ([]*appDto.AuditLogResponse, error)
	GetAuditLogsByDateRange(ctx context.Context, startDate, endDate int64) ([]*appDto.AuditLogResponse, error)
	DeleteOldAuditLogs(ctx context.Context, cutoffDate int64) error
}

package repository

import (
	"context"
	"fmt"
	"net"

	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
	"github.com/edynnt/veloras-api/internal/auth/domain/repository"
	"github.com/edynnt/veloras-api/internal/shared/gen"
	"github.com/edynnt/veloras-api/pkg/response/msg"
	"github.com/edynnt/veloras-api/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type auditRepository struct {
	db *gen.Queries
}

// CreateAuditLog implements repository.AuditRepository.
func (a *auditRepository) CreateAuditLog(ctx context.Context, log *entity.AuditLog) (string, error) {
	convertUserId, err := utils.ConvertUUID(log.UserID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	param := gen.CreateAuditLogParams{
		UserID:       convertUserId,
		Action:       log.Action,
		ResourceType: log.ResourceType,
		ResourceID:   pgtype.Text{String: log.ResourceID, Valid: true},
		Details:      []byte("{}"), // Convert map to JSON bytes
		IpAddress:    parseIPAddress(log.IPAddress),
		UserAgent:    pgtype.Text{String: log.UserAgent, Valid: true},
	}

	result, err := a.db.CreateAuditLog(ctx, param)
	if err != nil {
		return "", fmt.Errorf("%s: %w", msg.FailedToCreateAuditLog, err)
	}

	return result.ID.String(), nil
}

// GetAuditLogsByUser implements repository.AuditRepository.
func (a *auditRepository) GetAuditLogsByUser(ctx context.Context, userId string, limit int) ([]*entity.AuditLog, error) {
	convertUserId, err := utils.ConvertUUID(userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.UserIdInvalid, err)
	}

	param := gen.GetAuditLogsByUserParams{
		UserID: convertUserId,
		Limit:  int32(limit),
	}

	results, err := a.db.GetAuditLogsByUser(ctx, param)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGetAuditLogs, err)
	}

	var entityResults []*entity.AuditLog
	for _, result := range results {
		var entityResult entity.AuditLog
		if err := utils.SafeCopy(&entityResult, &result); err != nil {
			return nil, err
		}
		entityResults = append(entityResults, &entityResult)
	}

	return entityResults, nil
}

// GetAuditLogsByResource implements repository.AuditRepository.
func (a *auditRepository) GetAuditLogsByResource(ctx context.Context, resourceType string, resourceId string) ([]*entity.AuditLog, error) {
	param := gen.GetAuditLogsByResourceParams{
		ResourceType: resourceType,
		ResourceID:   pgtype.Text{String: resourceId, Valid: true},
	}

	results, err := a.db.GetAuditLogsByResource(ctx, param)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGetAuditLogs, err)
	}

	var entityResults []*entity.AuditLog
	for _, result := range results {
		var entityResult entity.AuditLog
		if err := utils.SafeCopy(&entityResult, &result); err != nil {
			return nil, err
		}
		entityResults = append(entityResults, &entityResult)
	}

	return entityResults, nil
}

// GetAuditLogsByAction implements repository.AuditRepository.
func (a *auditRepository) GetAuditLogsByAction(ctx context.Context, action string, limit int) ([]*entity.AuditLog, error) {
	param := gen.GetAuditLogsByActionParams{
		Action: action,
		Limit:  int32(limit),
	}

	results, err := a.db.GetAuditLogsByAction(ctx, param)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGetAuditLogs, err)
	}

	var entityResults []*entity.AuditLog
	for _, result := range results {
		var entityResult entity.AuditLog
		if err := utils.SafeCopy(&entityResult, &result); err != nil {
			return nil, err
		}
		entityResults = append(entityResults, &entityResult)
	}

	return entityResults, nil
}

// GetAuditLogsByDateRange implements repository.AuditRepository.
func (a *auditRepository) GetAuditLogsByDateRange(ctx context.Context, startDate, endDate int64) ([]*entity.AuditLog, error) {
	param := gen.GetAuditLogsByDateRangeParams{
		CreatedAt:   pgtype.Int8{Int64: startDate, Valid: true},
		CreatedAt_2: pgtype.Int8{Int64: endDate, Valid: true},
	}

	results, err := a.db.GetAuditLogsByDateRange(ctx, param)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGetAuditLogs, err)
	}

	var entityResults []*entity.AuditLog
	for _, result := range results {
		var entityResult entity.AuditLog
		if err := utils.SafeCopy(&entityResult, &result); err != nil {
			return nil, err
		}
		entityResults = append(entityResults, &entityResult)
	}

	return entityResults, nil
}

// DeleteOldAuditLogs implements repository.AuditRepository.
func (a *auditRepository) DeleteOldAuditLogs(ctx context.Context, cutoffDate int64) error {
	param := pgtype.Int8{Int64: cutoffDate, Valid: true}
	err := a.db.DeleteOldAuditLogs(ctx, param)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToDeleteOldAuditLogs, err)
	}

	return nil
}

// parseIPAddress parses IP address string to net.IP
func parseIPAddress(ipStr string) net.IP {
	if ipStr == "" {
		return nil
	}
	
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil
	}
	
	return ip
}

func NewAuditRepository(db *pgxpool.Pool) repository.AuditRepository {
	queries := gen.New(db)
	return &auditRepository{db: queries}
}

package service

import (
	"context"
	"fmt"

	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
	auditRepo "github.com/edynnt/veloras-api/internal/auth/domain/repository"
	"github.com/edynnt/veloras-api/pkg/response/msg"
)

type auditService struct {
	auditRepo auditRepo.AuditRepository
}

// CreateAuditLog implements AuditService.
func (a *auditService) CreateAuditLog(ctx context.Context, req appDto.CreateAuditLogRequest) (string, error) {
	auditLog := &entity.AuditLog{
		UserID:       req.UserID,
		Action:       req.Action,
		ResourceType: req.ResourceType,
		ResourceID:   req.ResourceID,
		Details:      req.Details,
		IPAddress:    req.IPAddress,
		UserAgent:    req.UserAgent,
	}

	auditLogId, err := a.auditRepo.CreateAuditLog(ctx, auditLog)
	if err != nil {
		return "", fmt.Errorf("%s: %w", msg.FailedToCreateAuditLog, err)
	}

	return auditLogId, nil
}

// GetAuditLogsByUser implements AuditService.
func (a *auditService) GetAuditLogsByUser(ctx context.Context, userId string, limit int) ([]*appDto.AuditLogResponse, error) {
	auditLogs, err := a.auditRepo.GetAuditLogsByUser(ctx, userId, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGetAuditLogs, err)
	}

	var responses []*appDto.AuditLogResponse
	for _, log := range auditLogs {
		responses = append(responses, &appDto.AuditLogResponse{
			ID:           log.ID,
			UserID:       log.UserID,
			Action:       log.Action,
			ResourceType: log.ResourceType,
			ResourceID:   log.ResourceID,
			Details:      log.Details,
			IPAddress:    log.IPAddress,
			UserAgent:    log.UserAgent,
			CreatedAt:    log.CreatedAt,
		})
	}

	return responses, nil
}

// GetAuditLogsByResource implements AuditService.
func (a *auditService) GetAuditLogsByResource(ctx context.Context, resourceType, resourceId string) ([]*appDto.AuditLogResponse, error) {
	auditLogs, err := a.auditRepo.GetAuditLogsByResource(ctx, resourceType, resourceId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGetAuditLogs, err)
	}

	var responses []*appDto.AuditLogResponse
	for _, log := range auditLogs {
		responses = append(responses, &appDto.AuditLogResponse{
			ID:           log.ID,
			UserID:       log.UserID,
			Action:       log.Action,
			ResourceType: log.ResourceType,
			ResourceID:   log.ResourceID,
			Details:      log.Details,
			IPAddress:    log.IPAddress,
			UserAgent:    log.UserAgent,
			CreatedAt:    log.CreatedAt,
		})
	}

	return responses, nil
}

// GetAuditLogsByAction implements AuditService.
func (a *auditService) GetAuditLogsByAction(ctx context.Context, action string, limit int) ([]*appDto.AuditLogResponse, error) {
	auditLogs, err := a.auditRepo.GetAuditLogsByAction(ctx, action, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGetAuditLogs, err)
	}

	var responses []*appDto.AuditLogResponse
	for _, log := range auditLogs {
		responses = append(responses, &appDto.AuditLogResponse{
			ID:           log.ID,
			UserID:       log.UserID,
			Action:       log.Action,
			ResourceType: log.ResourceType,
			ResourceID:   log.ResourceID,
			Details:      log.Details,
			IPAddress:    log.IPAddress,
			UserAgent:    log.UserAgent,
			CreatedAt:    log.CreatedAt,
		})
	}

	return responses, nil
}

// GetAuditLogsByDateRange implements AuditService.
func (a *auditService) GetAuditLogsByDateRange(ctx context.Context, startDate, endDate int64) ([]*appDto.AuditLogResponse, error) {
	auditLogs, err := a.auditRepo.GetAuditLogsByDateRange(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", msg.FailedToGetAuditLogs, err)
	}

	var responses []*appDto.AuditLogResponse
	for _, log := range auditLogs {
		responses = append(responses, &appDto.AuditLogResponse{
			ID:           log.ID,
			UserID:       log.UserID,
			Action:       log.Action,
			ResourceType: log.ResourceType,
			ResourceID:   log.ResourceID,
			Details:      log.Details,
			IPAddress:    log.IPAddress,
			UserAgent:    log.UserAgent,
			CreatedAt:    log.CreatedAt,
		})
	}

	return responses, nil
}

// DeleteOldAuditLogs implements AuditService.
func (a *auditService) DeleteOldAuditLogs(ctx context.Context, cutoffDate int64) error {
	err := a.auditRepo.DeleteOldAuditLogs(ctx, cutoffDate)
	if err != nil {
		return fmt.Errorf("%s: %w", msg.FailedToDeleteOldAuditLogs, err)
	}

	return nil
}

func NewAuditService(auditRepo auditRepo.AuditRepository) AuditService {
	return &auditService{
		auditRepo: auditRepo,
	}
}

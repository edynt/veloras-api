package repository

import (
	"context"

	"github.com/edynnt/veloras-api/internal/permission/domain/model/entity"
)

type PermissisonRepository interface {
	GetPermissions(ctx context.Context) (*entity.Permission, error)
}

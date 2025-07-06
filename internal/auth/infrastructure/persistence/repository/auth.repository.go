package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
	"github.com/edynnt/veloras-api/internal/auth/domain/repository"
)

type authRepository struct {
	db *sql.DB
}

// CreatUser implements repository.AuthRepository.
func (a *authRepository) CreatUser(ctx context.Context, account *entity.Account) (int64, error) {
	fmt.Println("call create account auth.repo infrastructure")
	panic("unimplemented")
}

func NewAuthRepository(db *sql.DB) repository.AuthRepository {
	return &authRepository{db}
}

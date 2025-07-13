package repository

import (
	"context"
	"fmt"

	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
	"github.com/edynnt/veloras-api/internal/auth/domain/repository"
	authsqlc "github.com/edynnt/veloras-api/internal/auth/infrastructure/persistence/sqlc/gen"
	"github.com/jackc/pgx/v5/pgxpool"
)

type authRepository struct {
	db *authsqlc.Queries
}

// CreatUser implements repository.AuthRepository.
func (a *authRepository) CreatUser(ctx context.Context, account *entity.Account) (int64, error) {
	fmt.Println("call create user auth.repository persistence")
	panic("unimplemented")
}

// UsernameExists implements repository.AuthRepository.
func (ar *authRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	fmt.Println("call username exists auth.repository persistence")
	panic("unimplemented")
}

func NewAuthRepository(db *pgxpool.Pool) repository.AuthRepository {
	queries := authsqlc.New(db) // db is *pgxpool.Pool
	return &authRepository{db: queries}
}

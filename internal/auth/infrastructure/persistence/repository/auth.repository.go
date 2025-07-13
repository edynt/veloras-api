package repository

import (
	"context"
	"fmt"

	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
	"github.com/edynnt/veloras-api/internal/auth/domain/repository"
	"github.com/edynnt/veloras-api/internal/auth/infrastructure/persistence/sqlc/gen"
	authsqlc "github.com/edynnt/veloras-api/internal/auth/infrastructure/persistence/sqlc/gen"
	"github.com/jackc/pgx/v5/pgxpool"
)

type authRepository struct {
	db *authsqlc.Queries
}

// EmailExists implements repository.AuthRepository.
func (a *authRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	return a.db.GetUsernameExists(ctx, email)
}

// CreateUser implements repository.AuthRepository.
func (a *authRepository) CreateUser(ctx context.Context, account *entity.Account) (int64, error) {
	createdAccount, err := a.db.CreateUser(ctx, gen.CreateUserParams{
		Username:    account.Username,
		Email:       account.Email,
		Password:    account.Password,
		PhoneNumber: account.PhoneNumber,
		FirstName:   account.FirstName,
		LastName:    account.LastName,
	})

	fmt.Println("createdAccount: ", createdAccount)
	fmt.Println("err: ", err)

	if err != nil {
		return 0, err
	}

	// if createdAccount == nil {
	// 	return 0, fmt.Errorf("account creation did not return an account instance")
	// }

	return 123, nil
}

// UsernameExists implements repository.AuthRepository.
func (a *authRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	return a.db.GetUsernameExists(ctx, username)
}

func NewAuthRepository(db *pgxpool.Pool) repository.AuthRepository {
	queries := authsqlc.New(db) // db is *pgxpool.Pool
	return &authRepository{db: queries}
}

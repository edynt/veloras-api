package repository

import (
	"context"

	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
	"github.com/edynnt/veloras-api/internal/auth/domain/repository"
	"github.com/edynnt/veloras-api/internal/shared/gen"
	authsqlc "github.com/edynnt/veloras-api/internal/shared/gen"
	"github.com/edynnt/veloras-api/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jinzhu/copier"
)

type authRepository struct {
	db *authsqlc.Queries
}

// CreateVerificationCode implements repository.AuthRepository.
func (a *authRepository) CreateVerificationCode(ctx context.Context, userVerification *entity.EmailVerification) (string, error) {
	var param gen.CreateEmailVerificationParams
	copier.Copy(&param, &userVerification)
	createdEmailVerification, err := a.db.CreateEmailVerification(ctx, param)

	if err != nil {
		return "", err
	}

	return utils.Int32ToString(createdEmailVerification.ID), nil

}

// EmailExists implements repository.AuthRepository.
func (a *authRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	return a.db.GetUsernameExists(ctx, email)
}

// CreateUser implements repository.AuthRepository.
func (a *authRepository) CreateUser(ctx context.Context, account *entity.Account) (string, error) {

	var param gen.CreateUserParams
	copier.Copy(&param, &account)

	createdAccount, err := a.db.CreateUser(ctx, param)

	if err != nil {
		return "", err
	}

	return createdAccount.ID.String(), nil
}

// UsernameExists implements repository.AuthRepository.
func (a *authRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	return a.db.GetUsernameExists(ctx, username)
}

func NewAuthRepository(db *pgxpool.Pool) repository.AuthRepository {
	queries := authsqlc.New(db) // db is *pgxpool.Pool
	return &authRepository{db: queries}
}

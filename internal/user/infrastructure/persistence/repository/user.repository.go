package repository

import (
	"context"
	"time"

	"github.com/edynnt/veloras-api/internal/shared/gen"
	"github.com/edynnt/veloras-api/internal/user/domain/model/entity"
	"github.com/edynnt/veloras-api/internal/user/domain/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) GetAllUsers(ctx context.Context, limit, offset int32) ([]entity.User, error) {
	queries := gen.New(r.db)
	
	rows, err := queries.GetAllUsers(ctx, gen.GetAllUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	users := make([]entity.User, 0, len(rows))
	for _, row := range rows {
		user := entity.User{
			ID:          row.ID,
			Email:       row.Email,
			Username:    row.Username,
			IsVerified:  row.IsVerified.Bool,
			PhoneNumber: row.PhoneNumber,
			FirstName:   row.FirstName,
			LastName:    row.LastName,
			Status:      row.Status.Int32,
			Language:    row.Language.String,
			CreatedAt:   time.Unix(row.CreatedAt.Int64, 0),
			UpdatedAt:   time.Unix(row.UpdatedAt.Int64, 0),
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepositoryImpl) CountAllUsers(ctx context.Context) (int64, error) {
	queries := gen.New(r.db)
	return queries.CountAllUsers(ctx)
}

func (r *userRepositoryImpl) GetUserByID(ctx context.Context, userID int32) (*entity.User, error) {
	queries := gen.New(r.db)
	
	row, err := queries.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		ID:          row.ID,
		Email:       row.Email,
		Username:    row.Username,
		IsVerified:  row.IsVerified.Bool,
		PhoneNumber: row.PhoneNumber,
		FirstName:   row.FirstName,
		LastName:    row.LastName,
		Status:      row.Status.Int32,
		Language:    row.Language.String,
		CreatedAt:   time.Unix(row.CreatedAt.Int64, 0),
		UpdatedAt:   time.Unix(row.UpdatedAt.Int64, 0),
	}

	return user, nil
}

func (r *userRepositoryImpl) UpdateUserProfile(ctx context.Context, userID int32, username, phoneNumber, firstName, lastName, language string) (*entity.User, error) {
	queries := gen.New(r.db)
	
	row, err := queries.UpdateUserProfile(ctx, gen.UpdateUserProfileParams{
		ID:      userID,
		Column2: username,
		Column3: phoneNumber,
		Column4: firstName,
		Column5: lastName,
		Column6: language,
	})
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		ID:          row.ID,
		Email:       row.Email,
		Username:    row.Username,
		PhoneNumber: row.PhoneNumber,
		FirstName:   row.FirstName,
		LastName:    row.LastName,
		Language:    row.Language.String,
		UpdatedAt:   time.Unix(row.UpdatedAt.Int64, 0),
	}

	return user, nil
}

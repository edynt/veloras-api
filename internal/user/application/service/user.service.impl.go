package service

import (
	"context"
	"errors"

	"github.com/edynt/chogiare/veloras-api/internal/user/application/service/dto"
	"github.com/edynt/chogiare/veloras-api/internal/user/domain/repository"
	"github.com/edynt/chogiare/veloras-api/pkg/utils"
)

type userServiceImpl struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userServiceImpl{
		userRepo: userRepo,
	}
}

func (s *userServiceImpl) GetAllUsers(ctx context.Context, page, pageSize int32) (*dto.UserListAppDTO, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	users, err := s.userRepo.GetAllUsers(ctx, pageSize, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.userRepo.CountAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	userDTOs := make([]dto.UserAppDTO, 0, len(users))
	for _, user := range users {
		userDTOs = append(userDTOs, dto.UserAppDTO{
			ID:          user.ID,
			Email:       user.Email,
			Username:    user.Username,
			IsVerified:  user.IsVerified,
			PhoneNumber: user.PhoneNumber,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Status:      user.Status,
			Language:    user.Language,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
		})
	}

	totalPages := int32((total + int64(pageSize) - 1) / int64(pageSize))

	return &dto.UserListAppDTO{
		Users:      userDTOs,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *userServiceImpl) GetUserByID(ctx context.Context, userID int32) (*dto.UserAppDTO, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserAppDTO{
		ID:          user.ID,
		Email:       user.Email,
		Username:    user.Username,
		IsVerified:  user.IsVerified,
		PhoneNumber: user.PhoneNumber,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Status:      user.Status,
		Language:    user.Language,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}

func (s *userServiceImpl) GetCurrentUser(ctx context.Context) (*dto.UserAppDTO, error) {
	subject := ctx.Value("subjectID")
	if subject == nil {
		return nil, errors.New("user not authenticated")
	}

	userID := utils.StringToInt(subject.(string))
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}

	return s.GetUserByID(ctx, int32(userID))
}

func (s *userServiceImpl) UpdateUserProfile(ctx context.Context, userID int32, req dto.UpdateUserProfileAppDTO) (*dto.UserAppDTO, error) {
	_, err := s.userRepo.UpdateUserProfile(ctx, userID, req.Username, req.PhoneNumber, req.FirstName, req.LastName, req.Language)
	if err != nil {
		return nil, err
	}

	// Get full user info after update
	return s.GetUserByID(ctx, userID)
}

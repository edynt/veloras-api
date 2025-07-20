package service

import (
	"context"
	"fmt"

	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
	"github.com/edynnt/veloras-api/internal/auth/domain/model/entity"
	authRepo "github.com/edynnt/veloras-api/internal/auth/domain/repository"
	"github.com/edynnt/veloras-api/utils"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	authRepo authRepo.AuthRepository
}

// Create implements AuthService.
func (as *authService) CreateUser(ctx context.Context, accountDto appDto.AccountAppDTO) (string, error) {
	fmt.Println("call create user auth.service.impl")

	//1. Check permissions -> event registered

	// 2. Check username exists
	exists, err := as.authRepo.UsernameExists(ctx, accountDto.Username)

	if err != nil {
		return "", fmt.Errorf("failed to check username exist: %w", err)
	}
	if exists {
		return "", fmt.Errorf("username already exists")
	}

	// 3. Check email exists
	exists, err = as.authRepo.EmailExists(ctx, accountDto.Email)
	if err != nil {
		return "", fmt.Errorf("failed to check email exist: %w", err)
	}

	if exists {
		return "", fmt.Errorf("email already exists")
	}

	// 4. GenerateFromPassword
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(accountDto.Password), bcrypt.DefaultCost)
	if err != nil {
		// log.Printf("Error hashing password for user %s: %v", accountDto.Username, err)
		return "", fmt.Errorf("failed to secure password: %w", err) // Không lộ chi tiết lỗi hash
	}
	hashedPassword := string(hashedPasswordBytes)

	accountDto.Password = hashedPassword

	// 5. Insert account into database
	newAccountId, err := as.authRepo.CreateUser(ctx, &entity.Account{
		Username:    accountDto.Username,
		Email:       accountDto.Email,
		Password:    accountDto.Password,
		PhoneNumber: accountDto.PhoneNumber,
		FirstName:   accountDto.FirstName,
		LastName:    accountDto.LastName,
	})

	if err != nil {
		// log.Printf("Error creating user %s in database: %v", accountToCreate.Username, err)
		return "", fmt.Errorf("could not create account: %w", err)
	}

	if newAccountId == "" {
		// log.Printf("Repository returned invalid account ID %d for user %s", newAccountId, accountToCreate.Username)
		return "", fmt.Errorf("account creation failed to return a valid ID")
	}

	as.authRepo.CreateVerificationCode(ctx, &entity.CreateVerificationCode{
		UserID:    newAccountId,
		Code:      utils.GenerateSixDigitCode(),
		ExpiresAt: utils.AddHours(1),
	})

	// 6. Return account ID
	return newAccountId, nil
}

func NewAuthService(
	authRepo authRepo.AuthRepository,
) AuthService {
	return &authService{
		authRepo: authRepo,
	}
}

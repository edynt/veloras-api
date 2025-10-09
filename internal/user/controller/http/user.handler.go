package http

import (
	"net/http"
	"strconv"

	"github.com/edynnt/veloras-api/internal/user/application/service"
	appDto "github.com/edynnt/veloras-api/internal/user/application/service/dto"
	ctlDto "github.com/edynnt/veloras-api/internal/user/controller/dto"
	"github.com/edynnt/veloras-api/pkg/response"
	"github.com/edynnt/veloras-api/pkg/response/msg"
	"github.com/edynnt/veloras-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GetAllUsers
// @Summary Get all users
// @Description Get a paginated list of all users
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} ctlDto.UserListResponse "Returns paginated list of users"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /users [get]
func (h *UserHandler) GetAllUsers(ctx *gin.Context) (res interface{}, err error) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 32)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	userList, err := h.service.GetAllUsers(ctx, int32(page), int32(pageSize))
	if err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, err.Error())
	}

	users := make([]ctlDto.UserResponse, 0, len(userList.Users))
	for _, user := range userList.Users {
		users = append(users, ctlDto.UserResponse{
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

	return ctlDto.UserListResponse{
		Users:      users,
		Total:      userList.Total,
		Page:       userList.Page,
		PageSize:   userList.PageSize,
		TotalPages: userList.TotalPages,
	}, nil
}

// GetUserByID
// @Summary Get user by ID
// @Description Get detailed information about a specific user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} ctlDto.UserResponse "Returns user information"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Failure 404 {object} response.APIError "User not found"
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(ctx *gin.Context) (res interface{}, err error) {
	userIDStr := ctx.Param("id")
	userID := utils.StringToInt(userIDStr)
	
	if userID == 0 {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, "Invalid user ID")
	}

	user, err := h.service.GetUserByID(ctx, int32(userID))
	if err != nil {
		return nil, response.NewAPIError(http.StatusNotFound, msg.InvalidRequest, err.Error())
	}

	return ctlDto.UserResponse{
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

// GetCurrentUser
// @Summary Get current user profile
// @Description Get the profile information of the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ctlDto.UserResponse "Returns current user information"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /users/me [get]
func (h *UserHandler) GetCurrentUser(ctx *gin.Context) (res interface{}, err error) {
	user, err := h.service.GetCurrentUser(ctx)
	if err != nil {
		return nil, response.NewAPIError(http.StatusUnauthorized, msg.Unauthorized, err.Error())
	}

	return ctlDto.UserResponse{
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

// UpdateUserProfile
// @Summary Update user profile
// @Description Update the profile information of the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ctlDto.UpdateUserProfileRequest true "Update profile request"
// @Success 200 {object} ctlDto.UserResponse "Returns updated user information"
// @Failure 400 {object} response.APIError "Invalid request"
// @Failure 401 {object} response.APIError "Unauthorized"
// @Router /users/me [put]
func (h *UserHandler) UpdateUserProfile(ctx *gin.Context) (res interface{}, err error) {
	var req ctlDto.UpdateUserProfileRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, err.Error())
	}

	validation, exists := ctx.Get("validation")
	if !exists {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, msg.ValidationNotFoundInContext)
	}

	if apiErr := utils.ValidateStruct(req, validation.(*validator.Validate)); apiErr != nil {
		return nil, apiErr
	}

	// Get user ID from context
	subject := ctx.Value("subjectID")
	if subject == nil {
		return nil, response.NewAPIError(http.StatusUnauthorized, msg.Unauthorized, "")
	}

	userID := utils.StringToInt(subject.(string))
	if userID == 0 {
		return nil, response.NewAPIError(http.StatusUnauthorized, msg.Unauthorized, msg.UserIdInvalid)
	}

	updateDTO := appDto.UpdateUserProfileAppDTO{
		Username:    req.Username,
		PhoneNumber: req.PhoneNumber,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Language:    req.Language,
	}

	user, err := h.service.UpdateUserProfile(ctx, int32(userID), updateDTO)
	if err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, err.Error())
	}

	return ctlDto.UserResponse{
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

package http

import (
	"net/http"

	"github.com/edynnt/veloras-api/internal/auth/application/service"
	"github.com/edynnt/veloras-api/internal/auth/application/service/dto"
	ctlDto "github.com/edynnt/veloras-api/internal/auth/controller/dto"
	"github.com/edynnt/veloras-api/pkg/response"
	"github.com/edynnt/veloras-api/pkg/response/msg"
	"github.com/edynnt/veloras-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type RoleHandler struct {
	service service.RoleService
}

// CreateRole
// @Summary Create a new role
// @Description Create a new role with name and description
// @Tags Roles
// @Accept json
// @Produce json
// @Param request body ctlDto.RoleReq true "Role creation request"
// @Success 200 {object} map[string]interface{} "Returns created role"
// @Failure 400 {object} response.APIError "Invalid request or validation errors"
// @Router /roles [post]
func (rh *RoleHandler) CreateRole(ctx *gin.Context) (res interface{}, err error) {
	var req ctlDto.RoleReq

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

	requestRole := dto.RoleAppDTO{
		Name:        req.Name,
		Description: req.Description,
	}

	res, err = rh.service.CreateRole(ctx, requestRole)

	if err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, err.Error())
	}

	return res, nil
}

// UpdateRole
// @Summary Update an existing role
// @Description Update a role's name and description by ID
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Param request body ctlDto.RoleReq true "Role update request"
// @Success 200 {object} map[string]interface{} "Returns updated role"
// @Failure 400 {object} response.APIError "Invalid request or validation errors"
// @Router /roles/{id} [put]
func (rh *RoleHandler) UpdateRole(ctx *gin.Context) (res interface{}, err error) {
	var req ctlDto.RoleReq

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

	requestRole := dto.RoleAppDTO{
		ID:          ctx.Param("id"),
		Name:        req.Name,
		Description: req.Description,
	}

	res, err = rh.service.UpdateRole(ctx, requestRole)

	if err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, err.Error())
	}

	return res, nil
}

// DeleteRole
// @Summary Delete a role
// @Description Delete a role by ID
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} map[string]interface{} "Role deleted successfully"
// @Failure 400 {object} response.APIError "Invalid request or role not found"
// @Router /roles/{id} [delete]
func (rh *RoleHandler) DeleteRole(ctx *gin.Context) (res interface{}, err error) {

	err = rh.service.DeleteRole(ctx, ctx.Param("id"))

	if err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, err.Error())
	}

	return res, nil
}

// GetRoles
// @Summary Get all roles
// @Description Retrieve a list of all available roles in the system
// @Tags Roles
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{} "Returns list of roles"
// @Failure 400 {object} response.APIError "Invalid request"
// @Router /roles [get]
func (rh *RoleHandler) GetRoles(ctx *gin.Context) (res interface{}, err error) {
	roles, err := rh.service.GetRoles(ctx)

	if err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, err.Error())
	}

	return roles, nil
}

// GetRole
// @Summary Get a role by ID
// @Description Retrieve a specific role by its ID
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} map[string]interface{} "Returns role details"
// @Failure 400 {object} response.APIError "Invalid request or role not found"
// @Router /roles/{id} [get]
func (rh *RoleHandler) GetRole(ctx *gin.Context) (res interface{}, err error) {
	role, err := rh.service.GetRoleById(ctx, ctx.Param("id"))

	if err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, err.Error())
	}

	return role, nil
}

func NewRoleHandler(service service.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

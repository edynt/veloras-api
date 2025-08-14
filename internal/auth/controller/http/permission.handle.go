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

type PermissionHandler struct {
	service service.PermissionService
}

func NewPermissionHandler(service service.PermissionService) *PermissionHandler {
	return &PermissionHandler{service: service}
}

// GetPermissions
// @Summary Get all permissions
// @Description Retrieve a list of all available permissions in the system
// @Tags Permissions
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{} "Returns list of permissions"
// @Failure 400 {object} response.APIError "Invalid request"
// @Router /permissions [get]
func (ph *PermissionHandler) GetPermissions(ctx *gin.Context) (res interface{}, err error) {
	permissions, err := ph.service.GetPermissions(ctx)

	if err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, err.Error())
	}

	return permissions, nil
}

// CreatePermission
// @Summary Create a new permission
// @Description Create a new permission with name and description
// @Tags Permissions
// @Accept json
// @Produce json
// @Param request body ctlDto.PermissionReq true "Permission creation request"
// @Success 200 {object} map[string]interface{} "Returns created permission"
// @Failure 400 {object} response.APIError "Invalid request or validation errors"
// @Router /permissions [post]
func (ph *PermissionHandler) CreatePermission(ctx *gin.Context) (res interface{}, err error) {
	var req ctlDto.PermissionReq

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

	requestPermission := dto.PermissionAppDTO{
		Name:        req.Name,
		Description: req.Description,
	}

	res, err = ph.service.CreatePermission(ctx, requestPermission)

	if err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, err.Error())
	}

	return res, nil
}

// UpdatePermission
// @Summary Update an existing permission
// @Description Update a permission's name and description by ID
// @Tags Permissions
// @Accept json
// @Produce json
// @Param id path string true "Permission ID"
// @Param request body ctlDto.PermissionReq true "Permission update request"
// @Success 200 {object} map[string]interface{} "Returns updated permission"
// @Failure 400 {object} response.APIError "Invalid request or validation errors"
// @Router /permissions/{id} [put]
func (ph *PermissionHandler) UpdatePermission(ctx *gin.Context) (res interface{}, err error) {
	var req ctlDto.PermissionReq

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

	requestPermission := dto.PermissionAppDTO{
		ID:          ctx.Param("id"),
		Name:        req.Name,
		Description: req.Description,
	}

	res, err = ph.service.UpdatePermission(ctx, requestPermission)

	if err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, err.Error())
	}

	return res, nil
}

// DeletePermission
// @Summary Delete a permission
// @Description Delete a permission by ID
// @Tags Permissions
// @Accept json
// @Produce json
// @Param id path string true "Permission ID"
// @Success 200 {object} map[string]interface{} "Permission deleted successfully"
// @Failure 400 {object} response.APIError "Invalid request or permission not found"
// @Router /permissions/{id} [delete]
func (ph *PermissionHandler) DeletePermission(ctx *gin.Context) (res interface{}, err error) {

	err = ph.service.DeletePermission(ctx, ctx.Param("id"))

	if err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, err.Error())
	}

	return res, nil
}

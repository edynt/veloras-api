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

func (ph *PermissionHandler) GetPermissions(ctx *gin.Context) (res interface{}, err error) {
	permissions, err := ph.service.GetPermissions(ctx)

	if err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, err.Error())
	}

	return permissions, nil
}

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

func (ph *PermissionHandler) DeletePermission(ctx *gin.Context) (res interface{}, err error) {

	err = ph.service.DeletePermission(ctx, ctx.Param("id"))

	if err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, err.Error())
	}

	return res, nil
}

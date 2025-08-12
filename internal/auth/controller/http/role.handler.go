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

func (rh *RoleHandler) DeleteRole(ctx *gin.Context) (res interface{}, err error) {

	err = rh.service.DeleteRole(ctx, ctx.Param("id"))

	if err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, err.Error())
	}

	return res, nil
}

func (rh *RoleHandler) GetRoles(ctx *gin.Context) (res interface{}, err error) {
	roles, err := rh.service.GetRoles(ctx)

	if err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, msg.InvalidRequest, err.Error())
	}

	return roles, nil
}

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

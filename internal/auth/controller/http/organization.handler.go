package http

import (
	"net/http"

	appDto "github.com/edynnt/veloras-api/internal/auth/application/service/dto"
	ctlDto "github.com/edynnt/veloras-api/internal/auth/controller/dto"
	"github.com/edynnt/veloras-api/internal/auth/application/service"
	"github.com/edynnt/veloras-api/pkg/response"
	"github.com/edynnt/veloras-api/pkg/response/msg"
	"github.com/edynnt/veloras-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type OrganizationHandler struct {
	orgService service.OrganizationService
}

// CreateOrganization
// @Summary Create a new organization
// @Description Create a new organization with the provided details
// @Tags Organizations
// @Accept json
// @Produce json
// @Param request body ctlDto.CreateOrganizationReq true "Organization creation request"
// @Success 201 {object} map[string]interface{} "Returns created organization ID"
// @Failure 400 {object} response.APIError "Invalid request format or validation errors"
// @Failure 500 {object} response.APIError "Internal server error"
// @Router /organizations [post]
func (oh *OrganizationHandler) CreateOrganization(ctx *gin.Context) (res interface{}, err error) {
	var req ctlDto.CreateOrganizationReq

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

	// TODO: Get user ID from JWT token context
	createdBy := "current-user-id" // Placeholder

	createReq := appDto.CreateOrganizationRequest{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
	}

	orgId, err := oh.orgService.CreateOrganization(ctx, createReq, createdBy)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, msg.FailedToCreateOrganization, err.Error())
	}

	return map[string]interface{}{
		"organization_id": orgId,
		"message":         "Organization created successfully",
	}, nil
}

// GetOrganizationById
// @Summary Get organization by ID
// @Description Retrieve organization details by organization ID
// @Tags Organizations
// @Accept json
// @Produce json
// @Param id path string true "Organization ID"
// @Success 200 {object} appDto.OrganizationResponse "Returns organization details"
// @Failure 400 {object} response.APIError "Invalid organization ID"
// @Failure 404 {object} response.APIError "Organization not found"
// @Router /organizations/{id} [get]
func (oh *OrganizationHandler) GetOrganizationById(ctx *gin.Context) (res interface{}, err error) {
	orgId := ctx.Param("id")

	org, err := oh.orgService.GetOrganizationById(ctx, orgId)
	if err != nil {
		return nil, response.NewAPIError(http.StatusNotFound, msg.FailedToGetOrganization, err.Error())
	}

	return org, nil
}

// GetOrganizationsByParent
// @Summary Get organizations by parent ID
// @Description Retrieve all organizations that have the specified parent organization
// @Tags Organizations
// @Accept json
// @Produce json
// @Param parentId path string true "Parent Organization ID"
// @Success 200 {array} appDto.OrganizationResponse "Returns list of organizations"
// @Failure 400 {object} response.APIError "Invalid parent organization ID"
// @Router /organizations/parent/{parentId} [get]
func (oh *OrganizationHandler) GetOrganizationsByParent(ctx *gin.Context) (res interface{}, err error) {
	parentId := ctx.Param("parentId")

	orgs, err := oh.orgService.GetOrganizationsByParent(ctx, parentId)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, msg.FailedToGetOrganization, err.Error())
	}

	return orgs, nil
}

// UpdateOrganization
// @Summary Update organization
// @Description Update organization details
// @Tags Organizations
// @Accept json
// @Produce json
// @Param id path string true "Organization ID"
// @Param request body ctlDto.UpdateOrganizationReq true "Organization update request"
// @Success 200 {object} map[string]interface{} "Returns success message"
// @Failure 400 {object} response.APIError "Invalid request format or validation errors"
// @Failure 404 {object} response.APIError "Organization not found"
// @Router /organizations/{id} [put]
func (oh *OrganizationHandler) UpdateOrganization(ctx *gin.Context) (res interface{}, err error) {
	var req ctlDto.UpdateOrganizationReq

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

	req.ID = ctx.Param("id")

	updateReq := appDto.UpdateOrganizationRequest{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
	}

	err = oh.orgService.UpdateOrganization(ctx, updateReq)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, msg.FailedToUpdateOrganization, err.Error())
	}

	return map[string]interface{}{
		"message": "Organization updated successfully",
	}, nil
}

// DeleteOrganization
// @Summary Delete organization
// @Description Delete an organization by ID
// @Tags Organizations
// @Accept json
// @Produce json
// @Param id path string true "Organization ID"
// @Success 200 {object} map[string]interface{} "Returns success message"
// @Failure 400 {object} response.APIError "Invalid organization ID"
// @Failure 404 {object} response.APIError "Organization not found"
// @Router /organizations/{id} [delete]
func (oh *OrganizationHandler) DeleteOrganization(ctx *gin.Context) (res interface{}, err error) {
	orgId := ctx.Param("id")

	err = oh.orgService.DeleteOrganization(ctx, orgId)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, msg.FailedToDeleteOrganization, err.Error())
	}

	return map[string]interface{}{
		"message": "Organization deleted successfully",
	}, nil
}

// AssignUserToOrganization
// @Summary Assign user to organization
// @Description Assign a user to a specific organization
// @Tags Organizations
// @Accept json
// @Produce json
// @Param request body ctlDto.AssignUserToOrganizationReq true "User assignment request"
// @Success 200 {object} map[string]interface{} "Returns success message"
// @Failure 400 {object} response.APIError "Invalid request format or validation errors"
// @Router /organizations/assign-user [post]
func (oh *OrganizationHandler) AssignUserToOrganization(ctx *gin.Context) (res interface{}, err error) {
	var req ctlDto.AssignUserToOrganizationReq

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

	err = oh.orgService.AssignUserToOrganization(ctx, req.UserID, req.OrganizationID)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, msg.FailedToAssignUserToOrganization, err.Error())
	}

	return map[string]interface{}{
		"message": "User assigned to organization successfully",
	}, nil
}

// GetOrganizationUsers
// @Summary Get organization users
// @Description Retrieve all users belonging to a specific organization
// @Tags Organizations
// @Accept json
// @Produce json
// @Param organizationId path string true "Organization ID"
// @Success 200 {array} appDto.UserResponse "Returns list of users"
// @Failure 400 {object} response.APIError "Invalid organization ID"
// @Router /organizations/{organizationId}/users [get]
func (oh *OrganizationHandler) GetOrganizationUsers(ctx *gin.Context) (res interface{}, err error) {
	organizationId := ctx.Param("organizationId")

	users, err := oh.orgService.GetOrganizationUsers(ctx, organizationId)
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, msg.FailedToGetOrganization, err.Error())
	}

	return users, nil
}

func NewOrganizationHandler(orgService service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{
		orgService: orgService,
	}
}

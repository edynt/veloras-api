package http

import (
	"fmt"

	"github.com/edynnt/veloras-api/internal/permission/application/service"
	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	service service.PermissionService
}

func NewPermissionHandler(service service.PermissionService) *PermissionHandler {
	return &PermissionHandler{service: service}
}

func (ph *PermissionHandler) GetPermissions(ctx *gin.Context) (res interface{}, err error) {
	a := 0
	fmt.Println(a)
	return nil, nil
}

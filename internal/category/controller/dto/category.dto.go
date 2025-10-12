package dto

import "time"

type CreateCategoryRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Slug        string  `json:"slug" validate:"required,min=1,max=255"`
	Description *string `json:"description,omitempty"`
	Image       *string `json:"image,omitempty"`
	ParentID    *string `json:"parentId,omitempty"`
	IsActive    bool    `json:"isActive"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Slug        *string `json:"slug,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty"`
	Image       *string `json:"image,omitempty"`
	ParentID    *string `json:"parentId,omitempty"`
	IsActive    *bool   `json:"isActive,omitempty"`
}

type CategoryResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  *string   `json:"description,omitempty"`
	Image        *string   `json:"image,omitempty"`
	ParentID     *string   `json:"parentId,omitempty"`
	ProductCount int32     `json:"productCount"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type CategoryListResponse struct {
	Categories []CategoryResponse `json:"categories"`
	Total      int64              `json:"total"`
	Page       int32              `json:"page"`
	PageSize   int32              `json:"pageSize"`
	TotalPages int32              `json:"totalPages"`
}

type CategoryStatsResponse struct {
	TotalCategories  int64 `json:"totalCategories"`
	ParentCategories int64 `json:"parentCategories"`
	SubCategories    int64 `json:"subCategories"`
	TotalProducts    int64 `json:"totalProducts"`
}

package dto

import "time"

type CreateProductRequest struct {
	Title         string   `json:"title" validate:"required,min=1,max=500"`
	Description   string   `json:"description" validate:"required,min=1"`
	Price         float64  `json:"price" validate:"required,min=0"`
	OriginalPrice *float64 `json:"originalPrice,omitempty" validate:"omitempty,min=0"`
	CategoryID    string   `json:"categoryId" validate:"required"`
	Images        []string `json:"images" validate:"required,min=1"`
	Condition     string   `json:"condition" validate:"required,oneof=new like_new good fair poor"`
	Tags          []string `json:"tags"`
	Location      string   `json:"location" validate:"required,min=1,max=255"`
	Stock         int32    `json:"stock" validate:"required,min=0"`
	Status        string   `json:"status" validate:"omitempty,oneof=draft active sold archived suspended"`
	Badges        []string `json:"badges"`
	IsFeatured    bool     `json:"isFeatured"`
	IsPromoted    bool     `json:"isPromoted"`
}

type UpdateProductRequest struct {
	Title         *string  `json:"title,omitempty" validate:"omitempty,min=1,max=500"`
	Description   *string  `json:"description,omitempty" validate:"omitempty,min=1"`
	Price         *float64 `json:"price,omitempty" validate:"omitempty,min=0"`
	OriginalPrice *float64 `json:"originalPrice,omitempty" validate:"omitempty,min=0"`
	CategoryID    *string  `json:"categoryId,omitempty"`
	Images        []string `json:"images,omitempty"`
	Condition     *string  `json:"condition,omitempty" validate:"omitempty,oneof=new like_new good fair poor"`
	Tags          []string `json:"tags,omitempty"`
	Location      *string  `json:"location,omitempty" validate:"omitempty,min=1,max=255"`
	Stock         *int32   `json:"stock,omitempty" validate:"omitempty,min=0"`
	Status        *string  `json:"status,omitempty" validate:"omitempty,oneof=draft active sold archived suspended"`
	Badges        []string `json:"badges,omitempty"`
	IsFeatured    *bool    `json:"isFeatured,omitempty"`
	IsPromoted    *bool    `json:"isPromoted,omitempty"`
}

type UpdateProductStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=draft active sold archived suspended"`
}

type UpdateProductStockRequest struct {
	Stock int32 `json:"stock" validate:"required,min=0"`
}

type BulkUpdateProductsRequest struct {
	IDs    []string             `json:"ids" validate:"required,min=1"`
	Update UpdateProductRequest `json:"update" validate:"required"`
}

type SearchProductsRequest struct {
	Query      *string  `json:"query,omitempty"`
	CategoryID *string  `json:"categoryId,omitempty"`
	MinPrice   *float64 `json:"minPrice,omitempty"`
	MaxPrice   *float64 `json:"maxPrice,omitempty"`
	Condition  *string  `json:"condition,omitempty"`
	Location   *string  `json:"location,omitempty"`
	SellerID   *int32   `json:"sellerId,omitempty"`
	Rating     *float64 `json:"rating,omitempty"`
	SortBy     *string  `json:"sortBy,omitempty"`
	SortOrder  *string  `json:"sortOrder,omitempty"`
	Page       *int32   `json:"page,omitempty"`
	Limit      *int32   `json:"limit,omitempty"`
}

type ProductResponse struct {
	ID            string            `json:"id"`
	Title         string            `json:"title"`
	Description   string            `json:"description"`
	Price         float64           `json:"price"`
	OriginalPrice *float64          `json:"originalPrice,omitempty"`
	CategoryID    string            `json:"categoryId"`
	Category      *CategoryResponse `json:"category,omitempty"`
	Images        []string          `json:"images"`
	Condition     string            `json:"condition"`
	Tags          []string          `json:"tags"`
	Location      string            `json:"location"`
	Stock         int32             `json:"stock"`
	SellerID      int32             `json:"sellerId"`
	Seller        *SellerResponse   `json:"seller,omitempty"`
	Store         *StoreResponse    `json:"store,omitempty"`
	Status        string            `json:"status"`
	Badges        []string          `json:"badges"`
	Rating        float64           `json:"rating"`
	ReviewCount   int32             `json:"reviewCount"`
	ViewCount     int32             `json:"viewCount"`
	IsFeatured    bool              `json:"isFeatured"`
	IsPromoted    bool              `json:"isPromoted"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
}

type CategoryResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Description  string `json:"description,omitempty"`
	Image        string `json:"image,omitempty"`
	ProductCount int32  `json:"productCount"`
	IsActive     bool   `json:"isActive"`
}

type SellerResponse struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type StoreResponse struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Logo   string  `json:"logo,omitempty"`
	Rating float64 `json:"rating"`
}

type ProductListResponse struct {
	Products   []ProductResponse `json:"products"`
	Total      int64             `json:"total"`
	Page       int32             `json:"page"`
	PageSize   int32             `json:"pageSize"`
	TotalPages int32             `json:"totalPages"`
}

type ProductStatsResponse struct {
	TotalProducts    int64    `json:"totalProducts"`
	ActiveProducts   int64    `json:"activeProducts"`
	DraftProducts    int64    `json:"draftProducts"`
	SoldProducts     int64    `json:"soldProducts"`
	FeaturedProducts int64    `json:"featuredProducts"`
	PromotedProducts int64    `json:"promotedProducts"`
	AverageRating    *float64 `json:"averageRating"`
	TotalViews       *int64   `json:"totalViews"`
}

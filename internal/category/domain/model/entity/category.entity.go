package entity

import (
	"time"

	"github.com/edynt/chogiare/veloras-api/internal/shared/gen"
	"github.com/jackc/pgx/v5/pgtype"
)

type Category struct {
	ID           string
	Name         string
	Slug         string
	Description  *string
	Image        *string
	ParentID     *string
	ProductCount int32
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CategoryWithChildren struct {
	Category
	Children []Category
}

type CategoryList struct {
	Categories []Category
	Total      int64
	Page       int32
	PageSize   int32
	TotalPages int32
}

type CategoryStats struct {
	TotalCategories  int64
	ParentCategories int64
	SubCategories    int64
	TotalProducts    int64
}

// Convert from SQLC model to domain entity
func FromSQLCCategory(c gen.Category) Category {
	return Category{
		ID:           c.ID.String(),
		Name:         c.Name,
		Slug:         c.Slug,
		Description:  convertTextPtr(c.Description),
		Image:        convertTextPtr(c.Image),
		ParentID:     convertUUIDPtr(c.ParentID),
		ProductCount: c.ProductCount.Int32,
		IsActive:     c.IsActive.Bool,
		CreatedAt:    c.CreatedAt.Time,
		UpdatedAt:    c.UpdatedAt.Time,
	}
}

// Helper function to convert pgtype.Text to *string
func convertTextPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

// Helper function to convert pgtype.UUID to *string
func convertUUIDPtr(u pgtype.UUID) *string {
	if !u.Valid {
		return nil
	}
	val := u.String()
	return &val
}

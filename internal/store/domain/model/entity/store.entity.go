package entity

import (
	"time"

	"github.com/edynt/chogiare/veloras-api/internal/shared/gen"
	"github.com/jackc/pgx/v5/pgtype"
)

type Store struct {
	ID            string
	UserID        int32
	Name          string
	Description   *string
	Logo          *string
	Banner        *string
	Website       *string
	Phone         *string
	Email         *string
	Address       *string
	City          *string
	State         *string
	Country       *string
	PostalCode    *string
	Rating        float64
	ReviewCount   int32
	ProductCount  int32
	FollowerCount int32
	IsVerified    bool
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type StoreWithDetails struct {
	Store
	UserName  *string
	UserEmail *string
}

type StoreStats struct {
	TotalStores    int64
	ActiveStores   int64
	VerifiedStores int64
	AverageRating  float64
	TotalProducts  int64
	TotalFollowers int64
}

// Convert from SQLC model to domain entity
func FromSQLCStore(s gen.Store) Store {
	return Store{
		ID:            s.ID.String(),
		UserID:        s.UserID,
		Name:          s.Name,
		Description:   convertTextPtr(s.Description),
		Logo:          convertTextPtr(s.Logo),
		Banner:        convertTextPtr(s.Banner),
		Website:       nil, // Not in generated model
		Phone:         convertTextPtr(s.Phone),
		Email:         convertTextPtr(s.Email),
		Address:       convertTextPtr(s.Address),
		City:          nil, // Not in generated model
		State:         nil, // Not in generated model
		Country:       nil, // Not in generated model
		PostalCode:    nil, // Not in generated model
		Rating:        convertNumericToFloat64(s.Rating),
		ReviewCount:   s.ReviewCount.Int32,
		ProductCount:  0, // Not in generated model
		FollowerCount: 0, // Not in generated model
		IsVerified:    s.IsVerified.Bool,
		IsActive:      true, // Default to active, not in generated model
		CreatedAt:     s.CreatedAt.Time,
		UpdatedAt:     s.UpdatedAt.Time,
	}
}

// Convert from SQLC model with details to domain entity
func FromSQLCStoreWithDetails(s gen.GetStoreRow) StoreWithDetails {
	return StoreWithDetails{
		Store: Store{
			ID:            s.ID.String(),
			UserID:        s.UserID,
			Name:          s.Name,
			Description:   convertTextPtr(s.Description),
			Logo:          convertTextPtr(s.Logo),
			Banner:        convertTextPtr(s.Banner),
			Website:       nil, // Not in generated model
			Phone:         convertTextPtr(s.Phone),
			Email:         convertTextPtr(s.Email),
			Address:       convertTextPtr(s.Address),
			City:          nil, // Not in generated model
			State:         nil, // Not in generated model
			Country:       nil, // Not in generated model
			PostalCode:    nil, // Not in generated model
			Rating:        convertNumericToFloat64(s.Rating),
			ReviewCount:   s.ReviewCount.Int32,
			ProductCount:  0, // Not in generated model
			FollowerCount: 0, // Not in generated model
			IsVerified:    s.IsVerified.Bool,
			IsActive:      true, // Default to active, not in generated model
			CreatedAt:     s.CreatedAt.Time,
			UpdatedAt:     s.UpdatedAt.Time,
		},
		UserName:  convertInterfacePtr(s.OwnerName), // Using OwnerName from generated model
		UserEmail: convertTextPtr(s.OwnerEmail),     // Using OwnerEmail from generated model
	}
}

// Helper function to convert pgtype.Numeric to float64
func convertNumericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0
	}
	val, _ := n.Float64Value()
	return val.Float64
}

// Helper function to convert pgtype.Text to *string
func convertTextPtr(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

// Helper function to convert interface{} to *string
func convertInterfacePtr(i interface{}) *string {
	if i == nil {
		return nil
	}
	if str, ok := i.(string); ok {
		return &str
	}
	return nil
}

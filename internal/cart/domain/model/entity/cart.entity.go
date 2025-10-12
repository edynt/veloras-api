package entity

import (
	"fmt"
	"time"

	"github.com/edynt/chogiare/veloras-api/internal/shared/gen"
	"github.com/jackc/pgx/v5/pgtype"
)

type Cart struct {
	ID        string
	UserID    int32
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CartItem struct {
	ID        string
	CartID    string
	ProductID string
	Quantity  int32
	Price     float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CartWithDetails struct {
	Cart
	Items []CartItemWithDetails
}

type CartItemWithDetails struct {
	CartItem
	ProductName   string
	ProductImage  string
	ProductPrice  float64
	ProductStock  int32
	ProductStatus string
}

type CartStats struct {
	TotalItems     int64
	TotalValue     float64
	UniqueProducts int64
}

// Convert from SQLC model to domain entity
func FromSQLCCart(userID int32) Cart {
	return Cart{
		ID:        fmt.Sprintf("cart-%d", userID), // Generate cart ID from user ID
		UserID:    userID,
		CreatedAt: time.Now(), // Use current time as cart doesn't have timestamps
		UpdatedAt: time.Now(),
	}
}

// Convert from SQLC model to domain entity
func FromSQLCCartItem(ci gen.CartItem) CartItem {
	return CartItem{
		ID:        ci.ID.String(),
		CartID:    fmt.Sprintf("cart-%d", ci.UserID), // Generate cart ID from user ID
		ProductID: ci.ProductID.String(),
		Quantity:  ci.Quantity,
		Price:     0.0, // Price not in generated model
		CreatedAt: ci.CreatedAt.Time,
		UpdatedAt: ci.UpdatedAt.Time,
	}
}

// Convert from SQLC model with details to domain entity
func FromSQLCCartWithDetails(userID int32, items []gen.ListCartItemsRow) CartWithDetails {
	cartItems := make([]CartItemWithDetails, len(items))
	for i, item := range items {
		cartItems[i] = CartItemWithDetails{
			CartItem: CartItem{
				ID:        item.ID.String(),
				CartID:    fmt.Sprintf("cart-%d", userID),
				ProductID: item.ProductID.String(),
				Quantity:  item.Quantity,
				Price:     convertNumericToFloat64(item.Price),
				CreatedAt: item.CreatedAt.Time,
				UpdatedAt: item.UpdatedAt.Time,
			},
			ProductName:   convertTextPtrToString(item.Title),
			ProductImage:  getFirstImage(item.Images),
			ProductPrice:  convertNumericToFloat64(item.Price),
			ProductStock:  int32(item.Stock.Int32),
			ProductStatus: convertTextPtrToString(item.Status),
		}
	}

	return CartWithDetails{
		Cart: Cart{
			ID:        fmt.Sprintf("cart-%d", userID),
			UserID:    userID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Items: cartItems,
	}
}

// Convert from SQLC model with details to domain entity
func FromSQLCCartItemWithDetails(ci gen.ListCartItemsRow, userID int32) CartItemWithDetails {
	return CartItemWithDetails{
		CartItem: CartItem{
			ID:        ci.ID.String(),
			CartID:    fmt.Sprintf("cart-%d", userID),
			ProductID: ci.ProductID.String(),
			Quantity:  ci.Quantity,
			Price:     convertNumericToFloat64(ci.Price),
			CreatedAt: ci.CreatedAt.Time,
			UpdatedAt: ci.UpdatedAt.Time,
		},
		ProductName:   convertTextPtrToString(ci.Title),
		ProductImage:  getFirstImage(ci.Images),
		ProductPrice:  convertNumericToFloat64(ci.Price),
		ProductStock:  int32(ci.Stock.Int32),
		ProductStatus: convertTextPtrToString(ci.Status),
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

// Helper function to convert pgtype.Text to string
func convertTextPtrToString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

// Convert from SQLC model with details to domain entity (for GetCartItemsWithDetailsRow)
func FromSQLCCartItemWithDetailsFromGet(ci gen.GetCartItemsWithDetailsRow, userID int32) CartItemWithDetails {
	return CartItemWithDetails{
		CartItem: CartItem{
			ID:        ci.ID.String(),
			CartID:    fmt.Sprintf("cart-%d", userID),
			ProductID: ci.ProductID.String(),
			Quantity:  ci.Quantity,
			Price:     convertNumericToFloat64(ci.Price),
			CreatedAt: ci.CreatedAt.Time,
			UpdatedAt: ci.UpdatedAt.Time,
		},
		ProductName:   convertTextPtrToString(ci.Title),
		ProductImage:  getFirstImage(ci.Images),
		ProductPrice:  convertNumericToFloat64(ci.Price),
		ProductStock:  int32(ci.Stock.Int32),
		ProductStatus: convertTextPtrToString(ci.Status),
	}
}

// Convert from SQLC model with details to domain entity (for GetCartItemWithDetailsRow)
func FromSQLCCartItemWithDetailsFromGetSingle(ci gen.GetCartItemWithDetailsRow, userID int32) CartItemWithDetails {
	return CartItemWithDetails{
		CartItem: CartItem{
			ID:        ci.ID.String(),
			CartID:    fmt.Sprintf("cart-%d", userID),
			ProductID: ci.ProductID.String(),
			Quantity:  ci.Quantity,
			Price:     convertNumericToFloat64(ci.Price),
			CreatedAt: ci.CreatedAt.Time,
			UpdatedAt: ci.UpdatedAt.Time,
		},
		ProductName:   convertTextPtrToString(ci.Title),
		ProductImage:  getFirstImage(ci.Images),
		ProductPrice:  convertNumericToFloat64(ci.Price),
		ProductStock:  int32(ci.Stock.Int32),
		ProductStatus: convertTextPtrToString(ci.Status),
	}
}

// Helper function to get first image from array
func getFirstImage(images []string) string {
	if len(images) > 0 {
		return images[0]
	}
	return ""
}

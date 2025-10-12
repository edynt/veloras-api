package entity

import (
	"fmt"
	"time"

	"github.com/edynt/chogiare/veloras-api/internal/shared/gen"
	"github.com/jackc/pgx/v5/pgtype"
)

type Order struct {
	ID              string
	UserID          int32
	StoreID         string
	Status          string
	PaymentStatus   string
	PaymentMethod   string
	Subtotal        float64
	Tax             float64
	Shipping        float64
	Discount        float64
	Total           float64
	Currency        string
	ShippingAddress string
	BillingAddress  string
	Notes           *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type OrderWithDetails struct {
	Order
	StoreName *string
	StoreLogo *string
	Items     []OrderItem
	UserEmail *string
	UserName  *string
}

type OrderItem struct {
	ID           string
	OrderID      string
	ProductID    string
	ProductName  string
	ProductImage string
	Price        float64
	Quantity     int32
	Subtotal     float64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type OrderStats struct {
	TotalOrders       int64
	PendingOrders     int64
	ProcessingOrders  int64
	ShippedOrders     int64
	DeliveredOrders   int64
	CancelledOrders   int64
	TotalRevenue      float64
	AverageOrderValue float64
}

// Convert from SQLC model to domain entity
func FromSQLCOrder(o gen.Order) Order {
	return Order{
		ID:              o.ID.String(),
		UserID:          o.BuyerID,                           // Using BuyerID from generated model
		StoreID:         fmt.Sprintf("store-%d", o.SellerID), // Generate store ID from seller ID
		Status:          o.Status,
		PaymentStatus:   o.PaymentStatus,
		PaymentMethod:   o.PaymentMethod,
		Subtotal:        convertNumericToFloat64(o.TotalAmount),
		Tax:             0.0, // Not in generated model
		Shipping:        0.0, // Not in generated model
		Discount:        0.0, // Not in generated model
		Total:           convertNumericToFloat64(o.TotalAmount),
		Currency:        "USD", // Default currency
		ShippingAddress: o.ShippingAddressID.String(),
		BillingAddress:  o.ShippingAddressID.String(), // Using same as shipping for now
		Notes:           convertTextPtr(o.Notes),
		CreatedAt:       o.CreatedAt.Time,
		UpdatedAt:       o.UpdatedAt.Time,
	}
}

// Convert from SQLC model with details to domain entity
func FromSQLCOrderWithDetails(o gen.GetOrderRow) OrderWithDetails {
	return OrderWithDetails{
		Order: Order{
			ID:              o.ID.String(),
			UserID:          o.BuyerID,                           // Using BuyerID from generated model
			StoreID:         fmt.Sprintf("store-%d", o.SellerID), // Generate store ID from seller ID
			Status:          o.Status,
			PaymentStatus:   o.PaymentStatus,
			PaymentMethod:   o.PaymentMethod,
			Subtotal:        convertNumericToFloat64(o.TotalAmount),
			Tax:             0.0, // Not in generated model
			Shipping:        0.0, // Not in generated model
			Discount:        0.0, // Not in generated model
			Total:           convertNumericToFloat64(o.TotalAmount),
			Currency:        "USD", // Default currency
			ShippingAddress: o.ShippingAddressID.String(),
			BillingAddress:  o.ShippingAddressID.String(), // Using same as shipping for now
			Notes:           convertTextPtr(o.Notes),
			CreatedAt:       o.CreatedAt.Time,
			UpdatedAt:       o.UpdatedAt.Time,
		},
		StoreName: nil,                              // Not in generated model
		StoreLogo: nil,                              // Not in generated model
		UserEmail: convertTextPtr(o.BuyerEmail),     // Using BuyerEmail from generated model
		UserName:  convertInterfacePtr(o.BuyerName), // Using BuyerName from generated model
	}
}

// OrderItem conversion not needed as the generated Order model includes product and quantity directly

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

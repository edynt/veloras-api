package entity

import (
	"time"

	"github.com/edynt/chogiare/veloras-api/internal/shared/gen"
	"github.com/jackc/pgx/v5/pgtype"
)

type Order struct {
	ID                string
	BuyerID           int32
	SellerID          int32
	ProductID         string
	Quantity          int32
	TotalAmount       float64
	Status            string
	PaymentMethod     string
	PaymentStatus     string
	ShippingAddressID string
	Notes             *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
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
		ID:                o.ID.String(),
		BuyerID:           o.BuyerID,
		SellerID:          o.SellerID,
		ProductID:         o.ProductID.String(),
		Quantity:          o.Quantity,
		TotalAmount:       convertNumericToFloat64(o.TotalAmount),
		Status:            o.Status,
		PaymentMethod:     o.PaymentMethod,
		PaymentStatus:     o.PaymentStatus,
		ShippingAddressID: o.ShippingAddressID.String(),
		Notes:             convertTextPtr(o.Notes),
		CreatedAt:         o.CreatedAt.Time,
		UpdatedAt:         o.UpdatedAt.Time,
	}
}

// Convert from SQLC model with details to domain entity
func FromSQLCOrderWithDetails(o gen.GetOrderRow) OrderWithDetails {
	return OrderWithDetails{
		Order: Order{
			ID:                o.ID.String(),
			BuyerID:           o.BuyerID,
			SellerID:          o.SellerID,
			ProductID:         o.ProductID.String(),
			Quantity:          o.Quantity,
			TotalAmount:       convertNumericToFloat64(o.TotalAmount),
			Status:            o.Status,
			PaymentMethod:     o.PaymentMethod,
			PaymentStatus:     o.PaymentStatus,
			ShippingAddressID: o.ShippingAddressID.String(),
			Notes:             convertTextPtr(o.Notes),
			CreatedAt:         o.CreatedAt.Time,
			UpdatedAt:         o.UpdatedAt.Time,
		},
		StoreName: nil,                              // Not in generated model
		StoreLogo: nil,                              // Not in generated model
		UserEmail: nil,                              // BuyerEmail not available in current schema     // Using BuyerEmail from generated model
		UserName:  convertInterfacePtr(o.BuyerName), // Using BuyerName from generated model
	}
}

// Convert from SQLC ListOrdersRow to domain entity
func FromSQLCOrderWithDetailsFromList(o gen.ListOrdersRow) OrderWithDetails {
	return OrderWithDetails{
		Order: Order{
			ID:                o.ID.String(),
			BuyerID:           o.BuyerID,
			SellerID:          o.SellerID,
			ProductID:         o.ProductID.String(),
			Quantity:          o.Quantity,
			TotalAmount:       convertNumericToFloat64(o.TotalAmount),
			Status:            o.Status,
			PaymentMethod:     o.PaymentMethod,
			PaymentStatus:     o.PaymentStatus,
			ShippingAddressID: o.ShippingAddressID.String(),
			Notes:             convertTextPtr(o.Notes),
			CreatedAt:         o.CreatedAt.Time,
			UpdatedAt:         o.UpdatedAt.Time,
		},
		StoreName: nil,
		StoreLogo: nil,
		UserEmail: nil, // BuyerEmail not available in current schema
		UserName:  convertInterfacePtr(o.BuyerName),
	}
}

// Convert from SQLC ListOrdersByUserRow to domain entity
func FromSQLCOrderWithDetailsFromListByUser(o gen.ListOrdersByUserRow) OrderWithDetails {
	return OrderWithDetails{
		Order: Order{
			ID:                o.ID.String(),
			BuyerID:           o.BuyerID,
			SellerID:          o.SellerID,
			ProductID:         o.ProductID.String(),
			Quantity:          o.Quantity,
			TotalAmount:       convertNumericToFloat64(o.TotalAmount),
			Status:            o.Status,
			PaymentMethod:     o.PaymentMethod,
			PaymentStatus:     o.PaymentStatus,
			ShippingAddressID: o.ShippingAddressID.String(),
			Notes:             convertTextPtr(o.Notes),
			CreatedAt:         o.CreatedAt.Time,
			UpdatedAt:         o.UpdatedAt.Time,
		},
		StoreName: nil,
		StoreLogo: nil,
		UserEmail: nil, // BuyerEmail not available in current schema
		UserName:  convertInterfacePtr(o.BuyerName),
	}
}

// Convert from SQLC ListOrdersByStoreRow to domain entity
func FromSQLCOrderWithDetailsFromListByStore(o gen.ListOrdersByStoreRow) OrderWithDetails {
	return OrderWithDetails{
		Order: Order{
			ID:                o.ID.String(),
			BuyerID:           o.BuyerID,
			SellerID:          o.SellerID,
			ProductID:         o.ProductID.String(),
			Quantity:          o.Quantity,
			TotalAmount:       convertNumericToFloat64(o.TotalAmount),
			Status:            o.Status,
			PaymentMethod:     o.PaymentMethod,
			PaymentStatus:     o.PaymentStatus,
			ShippingAddressID: o.ShippingAddressID.String(),
			Notes:             convertTextPtr(o.Notes),
			CreatedAt:         o.CreatedAt.Time,
			UpdatedAt:         o.UpdatedAt.Time,
		},
		StoreName: nil,
		StoreLogo: nil,
		UserEmail: nil, // BuyerEmail not available in current schema
		UserName:  convertInterfacePtr(o.BuyerName),
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

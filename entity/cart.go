package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Cart struct {
	ID          int
	CustomerID  int
	Status      CartStatus
	TotalAmount decimal.Decimal
	CartItems   []*CartItem `gorm:"-"` // ignore this field when write and read with struct
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   *time.Time
	UpdatedBy   *string
}

type CartStatus string

var (
	CartStatusOpen   CartStatus = "OPEN"
	CartStatusClosed CartStatus = "CLOSED"
)

type CartItem struct {
	ID        int
	CartID    int
	PlantID   int
	Quantity  int
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt *time.Time
	UpdatedBy *string
}

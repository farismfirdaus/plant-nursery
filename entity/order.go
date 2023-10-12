package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID          int
	CustomerID  int
	TotalAmount decimal.Decimal
	OrderItems  []*OrderItem `gorm:"-"` // ignore this field when write and read with struct
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   *time.Time
	UpdatedBy   *string
}

type OrderItem struct {
	ID          int
	OrderID     int
	PlantID     int
	Price       decimal.Decimal
	Quantity    int
	TotalAmount decimal.Decimal
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   *time.Time
	UpdatedBy   *string
}

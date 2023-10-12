package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Plant struct {
	ID        int
	Sku       string
	Name      string
	ImageUrl  string
	Stock     int
	Price     decimal.Decimal
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt *time.Time
	UpdatedBy *string
}

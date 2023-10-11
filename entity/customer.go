package entity

import "time"

type Customer struct {
	ID        int
	FullName  string
	Email     string
	Password  string
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt *time.Time
	UpdatedBy *string
}

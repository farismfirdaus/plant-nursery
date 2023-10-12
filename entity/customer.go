package entity

import (
	"time"

	apperr "github.com/farismfirdaus/plant-nursery/errors"
	"github.com/farismfirdaus/plant-nursery/utils/helper"
)

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

func (c *Customer) Validate() error {
	switch {
	case c == nil:
		return apperr.BadRequest
	case c.FullName == "":
		return apperr.InvalidCustomerFullName
	case !helper.ValidateEmailAddr(c.Email):
		return apperr.InvalidEmail
	case c.Password == "":
		return apperr.InvalidPassword
	}

	return nil
}

package main

import (
	repository_customer "github.com/farismfirdaus/plant-nursery/services/customer/repository"

	"gorm.io/gorm"
)

type repository struct {
	customerRepo repository_customer.Customer
}

func setupRepository(db *gorm.DB) *repository {
	return &repository{
		customerRepo: repository_customer.NewRepository(db),
	}
}

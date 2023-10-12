package main

import (
	"gorm.io/gorm"

	repository_customer "github.com/farismfirdaus/plant-nursery/services/customer/repository"
	repository_plant "github.com/farismfirdaus/plant-nursery/services/plant/repository"
)

type repository struct {
	customerRepo repository_customer.Customer
	plantRepo    repository_plant.Plant
}

func setupRepository(db *gorm.DB) *repository {
	return &repository{
		customerRepo: repository_customer.NewRepository(db),
		plantRepo:    repository_plant.NewRepository(db),
	}
}

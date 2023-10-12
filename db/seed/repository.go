package main

import (
	"gorm.io/gorm"

	repository_cart "github.com/farismfirdaus/plant-nursery/services/cart/repository"
	repository_customer "github.com/farismfirdaus/plant-nursery/services/customer/repository"
	repository_order "github.com/farismfirdaus/plant-nursery/services/order/repository"
	repository_plant "github.com/farismfirdaus/plant-nursery/services/plant/repository"
)

type repository struct {
	customerRepo repository_customer.Customer
	plantRepo    repository_plant.Plant
	cartRepo     repository_cart.Cart
	orderRepo    repository_order.Order
}

func setupRepository(db *gorm.DB) *repository {
	return &repository{
		customerRepo: repository_customer.NewRepository(db),
		plantRepo:    repository_plant.NewRepository(db),
		cartRepo:     repository_cart.NewRepository(db),
		orderRepo:    repository_order.NewRepository(db),
	}
}

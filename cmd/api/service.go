package main

import (
	"github.com/farismfirdaus/plant-nursery/auth"
	"github.com/farismfirdaus/plant-nursery/services/cart"
	"github.com/farismfirdaus/plant-nursery/services/customer"
	"github.com/farismfirdaus/plant-nursery/services/plant"
)

type service struct {
	customerService customer.Customer
	plantService    plant.Plant
	cartService     cart.Cart
}

func setupService(repo *repository, client auth.Auth) *service {
	customerSvc := customer.NewHandler(repo.customerRepo, client)
	plantSvc := plant.NewHandler(repo.plantRepo)
	cartSvc := cart.NewHandler(repo.cartRepo, plantSvc)

	return &service{
		customerService: customerSvc,
		plantService:    plantSvc,
		cartService:     cartSvc,
	}
}

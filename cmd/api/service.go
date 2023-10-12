package main

import (
	"github.com/farismfirdaus/plant-nursery/auth"
	"github.com/farismfirdaus/plant-nursery/services/customer"
	"github.com/farismfirdaus/plant-nursery/services/plant"
)

type service struct {
	customerService customer.Customer
	plantService    plant.Plant
}

func setupService(repo *repository, client auth.Auth) *service {
	return &service{
		customerService: customer.NewHandler(repo.customerRepo, client),
		plantService:    plant.NewHandler(repo.plantRepo),
	}
}

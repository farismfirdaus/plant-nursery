package main

import (
	"github.com/farismfirdaus/plant-nursery/auth"
	"github.com/farismfirdaus/plant-nursery/services/customer"
)

type service struct {
	customerService customer.Customer
}

func setupService(repo *repository, client auth.Auth) *service {
	return &service{
		customerService: customer.NewHandler(repo.customerRepo, client),
	}
}

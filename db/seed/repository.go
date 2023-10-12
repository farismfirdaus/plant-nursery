package main

import (
	repository_plant "github.com/farismfirdaus/plant-nursery/services/plant/repository"

	"gorm.io/gorm"
)

type repository struct {
	plantRepo repository_plant.Repository
}

func setupRepository(db *gorm.DB) *repository {
	return &repository{
		plantRepo: *repository_plant.NewRepository(db),
	}
}

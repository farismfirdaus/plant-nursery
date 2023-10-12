package main

import (
	"github.com/gin-gonic/gin"

	controller_customer "github.com/farismfirdaus/plant-nursery/services/customer/controller"
	controller_plant "github.com/farismfirdaus/plant-nursery/services/plant/controller"
)

func setupController(g *gin.RouterGroup, serv *service) {
	controller_customer.New(g, serv.customerService)
	controller_plant.New(g, serv.plantService)
}

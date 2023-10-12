package main

import (
	"github.com/gin-gonic/gin"

	controller_customer "github.com/farismfirdaus/plant-nursery/services/customer/controller"
)

func setupController(g *gin.RouterGroup, serv *service) {
	controller_customer.New(g, serv.customerService)
}

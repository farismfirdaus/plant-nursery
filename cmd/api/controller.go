package main

import (
	"github.com/gin-gonic/gin"

	controller_cart "github.com/farismfirdaus/plant-nursery/services/cart/controller"
	controller_customer "github.com/farismfirdaus/plant-nursery/services/customer/controller"
	controller_order "github.com/farismfirdaus/plant-nursery/services/order/controller"
	controller_plant "github.com/farismfirdaus/plant-nursery/services/plant/controller"
)

func setupController(g *gin.RouterGroup, serv *service) {
	controller_customer.New(g, serv.customerService)
	controller_plant.New(g, serv.plantService)
	controller_cart.New(g, serv.cartService)
	controller_order.New(g, serv.orderService)
}

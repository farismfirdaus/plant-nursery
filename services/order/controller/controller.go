package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	apperr "github.com/farismfirdaus/plant-nursery/errors"
	"github.com/farismfirdaus/plant-nursery/services/order"
	"github.com/farismfirdaus/plant-nursery/utils/middleware"
	"github.com/farismfirdaus/plant-nursery/utils/response"
)

type Controller struct {
	service order.Order
}

func New(g *gin.RouterGroup, service order.Order) {
	c := Controller{
		service: service,
	}
	g.GET("/orders", middleware.Authenticate(), c.GetList)
	g.GET("/orders/items", middleware.Authenticate(), c.GetListUniqueItems)
	g.POST("/orders", middleware.Authenticate(), c.Checkout)
}

func (c *Controller) GetList(ctx *gin.Context) {
	req, exists := ctx.Get(gin.AuthUserKey)
	if !exists {
		response.BuildErrors(ctx, apperr.Unauthorized)
		return
	}

	customerID, ok := req.(int)
	if !ok {
		response.BuildErrors(ctx, apperr.Unauthorized)
		return
	}

	orders, err := c.service.GetList(ctx.Request.Context(), customerID)
	if err != nil {
		response.BuildErrors(ctx, err)
		return
	}

	response.BuildSuccess(ctx, http.StatusOK, orders)
}

func (c *Controller) GetListUniqueItems(ctx *gin.Context) {
	req, exists := ctx.Get(gin.AuthUserKey)
	if !exists {
		response.BuildErrors(ctx, apperr.Unauthorized)
		return
	}

	customerID, ok := req.(int)
	if !ok {
		response.BuildErrors(ctx, apperr.Unauthorized)
		return
	}

	orderItems, err := c.service.GetListUniqueItems(ctx.Request.Context(), customerID)
	if err != nil {
		response.BuildErrors(ctx, err)
		return
	}

	response.BuildSuccess(ctx, http.StatusOK, orderItems)
}

func (c *Controller) Checkout(ctx *gin.Context) {
	req, exists := ctx.Get(gin.AuthUserKey)
	if !exists {
		response.BuildErrors(ctx, apperr.Unauthorized)
		return
	}

	customerID, ok := req.(int)
	if !ok {
		response.BuildErrors(ctx, apperr.Unauthorized)
		return
	}

	request := map[string]int{}
	if err := json.NewDecoder(ctx.Request.Body).Decode(&request); err != nil {
		response.BuildErrors(ctx, apperr.BadRequest)
		return
	}

	err := c.service.Checkout(ctx.Request.Context(), customerID, request["cartId"])
	if err != nil {
		response.BuildErrors(ctx, err)
		return
	}

	response.BuildSuccess(ctx, http.StatusCreated, "created")
}

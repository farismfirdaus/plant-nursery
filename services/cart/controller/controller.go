package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	apperr "github.com/farismfirdaus/plant-nursery/errors"
	"github.com/farismfirdaus/plant-nursery/services/cart"
	"github.com/farismfirdaus/plant-nursery/utils/middleware"
	"github.com/farismfirdaus/plant-nursery/utils/response"
)

type Controller struct {
	service cart.Cart
}

func New(g *gin.RouterGroup, service cart.Cart) {
	c := Controller{
		service: service,
	}
	g.GET("/carts", middleware.Authenticate(), c.Get)
	g.POST("/carts/items", middleware.Authenticate(), c.AddItems)
}

func (c *Controller) Get(ctx *gin.Context) {
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

	cart, err := c.service.Get(ctx.Request.Context(), customerID)
	if err != nil {
		response.BuildErrors(ctx, err)
		return
	}

	response.BuildSuccess(ctx, http.StatusOK, cart)
}

func (c *Controller) AddItems(ctx *gin.Context) {
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

	request := []*cart.AddItemsRequest{}
	if err := json.NewDecoder(ctx.Request.Body).Decode(&request); err != nil {
		response.BuildErrors(ctx, apperr.BadRequest)
		return
	}

	errs := c.service.AddItems(ctx.Request.Context(), customerID, request)
	if len(errs) > 0 {
		response.BuildErrors(ctx, errs...)
		return
	}

	response.BuildSuccess(ctx, http.StatusCreated, "created")
}

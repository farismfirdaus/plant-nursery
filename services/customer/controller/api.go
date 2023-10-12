package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/farismfirdaus/plant-nursery/entity"
	apperr "github.com/farismfirdaus/plant-nursery/errors"
	"github.com/farismfirdaus/plant-nursery/services/customer"
	"github.com/farismfirdaus/plant-nursery/utils/response"
)

type Controller struct {
	service customer.Customer
}

func New(g *gin.RouterGroup, service customer.Customer) {
	c := Controller{
		service: service,
	}
	g.POST("/customers", c.Register)
	g.POST("/sessions", c.NewSession)
}

func (c *Controller) Register(ctx *gin.Context) {
	customer := &entity.Customer{}
	if err := json.NewDecoder(ctx.Request.Body).Decode(&customer); err != nil {
		response.BuildErrors(ctx, apperr.BadRequest)
		return
	}

	if err := c.service.Register(ctx.Request.Context(), customer); err != nil {
		response.BuildErrors(ctx, err)
		return
	}

	response.BuildSuccess(ctx, http.StatusCreated, "created")
}

func (c *Controller) NewSession(ctx *gin.Context) {
	customer := &entity.Customer{} // TODO: should use different entity that hold only email and password
	if err := json.NewDecoder(ctx.Request.Body).Decode(&customer); err != nil {
		response.BuildErrors(ctx, apperr.BadRequest)
		return
	}

	token, err := c.service.NewSession(ctx.Request.Context(), customer.Email, customer.Password)
	if err != nil {
		response.BuildErrors(ctx, err)
		return
	}

	response.BuildSuccess(ctx, http.StatusOK, map[string]string{"token": token})
}

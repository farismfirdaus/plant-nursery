package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/farismfirdaus/plant-nursery/services/plant"
	"github.com/farismfirdaus/plant-nursery/utils/response"
)

type Controller struct {
	service plant.Plant
}

func New(g *gin.RouterGroup, service plant.Plant) {
	c := Controller{
		service: service,
	}
	g.GET("/plants", c.GetList)
}

func (c *Controller) GetList(ctx *gin.Context) {
	plants, err := c.service.GetList(ctx.Request.Context())
	if err != nil {
		response.BuildErrors(ctx, err)
		return
	}

	response.BuildSuccess(ctx, http.StatusOK, plants)
}

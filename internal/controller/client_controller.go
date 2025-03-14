package controller

import (
	"core/internal/api"
	"core/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ClientController interface {
	Signup(c *gin.Context)
}

type clientController struct {
	service service.ClientService
}

func (controller clientController) Signup(c *gin.Context) {
	request := &api.ClientRegister{}
	if err := c.ShouldBind(request); err != nil && errors.As(err, &validator.ValidationErrors{}) {
		api.GetErrorJSON(c, "JSON is invalid")
		return
	}
	err := controller.service.Signup(request)
	if err != nil {
		api.GetErrorJSON(c, "Internal error")
		return
	}

}

func NewClientController(service service.ClientService) ClientController {
	return &clientController{
		service: service,
	}
}

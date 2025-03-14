package controller

import (
	"core/internal"
	"core/internal/api"
	"core/internal/service"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
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
		api.GetErrorJSON(c, http.StatusBadRequest, "JSON is invalid")
		return
	}
	client, err := controller.service.Signup(request)
	if err != nil {
		api.GetErrorJSON(c, http.StatusPreconditionFailed, err.Error())
		return
	}
	c.JSON(http.StatusOK, api.ResponseClientRegister{
		StatusResponse: &internal.StatusResponse{Status: "success"},
		ResponseUser: &api.ResponseUser{
			ID:    client.ID,
			Token: "mock_token",
			Type:  client.Type,
		},
	})
}

func NewClientController(service service.ClientService) ClientController {
	return &clientController{
		service: service,
	}
}

package api

import (
	"core/internal"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetErrorJSON(c *gin.Context, description string) {
	response := internal.InfoResponse{
		StatusResponse: &internal.StatusResponse{Status: "error"},
		Description:    description,
	}
	c.JSON(http.StatusOK, response)
}

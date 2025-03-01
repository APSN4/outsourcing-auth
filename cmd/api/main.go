package main

import (
	"core/internal/api"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	v1 := r.Group("v1")
	{
		userGroup := r.Group("user")
		{
			userGroup.GET("/register", api.RegisterUser)
		}
	}

	r.Run()
}

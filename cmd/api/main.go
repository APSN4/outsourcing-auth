package main

import (
	"core/internal/api"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	v1 := r.Group("v1")
	{
		v1.GET("/login", api.GetMockStandard)
		v1.GET("/account", api.GetMockStandard)
		registerGroup := v1.Group("register")
		{
			registerGroup.GET("/client", api.GetMockStandard)
			registerGroup.GET("/company", api.GetMockStandard)
		}
	}

	r.Run()
}

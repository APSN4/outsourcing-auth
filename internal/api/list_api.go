package api

import "github.com/gin-gonic/gin"

func GetListAPI(c *gin.Context) {
	c.JSON(200, gin.H{
		"api": "test",
	})
}

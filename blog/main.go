package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "dasdfas",
		})
	})

	r.GET("/get", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "gget",
		})
	})

	r.GET("/asd", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "gget",
		})
	})

	r.Run(":8080")
}

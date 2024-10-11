package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("layout/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "메인페이지",
		})
	})
	r.POST("/login", Login)

	r.Run(":8080")
}

func Login(c *gin.Context) {

}

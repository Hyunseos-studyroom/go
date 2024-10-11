package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.POST("/login", Login)

	r.Run(":8080")
}

func Login(c *gin.Context) {
	var requestBody struct {
		Code string `json:"code"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
		return
	}

	// 여기에 실제 토큰 발급 로직 추가
	accessToken := "your_access_token"
	refreshToken := "your_refresh_token"

	// 쿠키 설정
	c.SetCookie("access_token", accessToken, 3600, "/", "localhost", false, true)         // Secure: true로 설정하면 HTTPS에서만 쿠키 전송
	c.SetCookie("refresh_token", refreshToken, 3600*24*30, "/", "localhost", false, true) // 30일 유효

	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
}

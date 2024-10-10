package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"log"
	"net/http"
	"os"
	"time"
)

var client *redis.Client

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

var user = User{
	ID:       1,
	Username: "username",
	Password: "password",
	Phone:    "12038485839",
}

var (
	router = gin.Default()
)

func main() {
	router.POST("/login", Login)
	log.Fatal(router.Run(":8080"))
}

func Login(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	if user.Username != u.Username || user.Password != u.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}
	token, err := CreateToken(user.ID)
	log.Println(token)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	c.JSON(http.StatusOK, token)
}

func CreateToken(userId uint64) (string, error) {
	var err error

	os.Setenv("ACCESS_SECRET", "secret key!")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

func init() {
	dsn := "localhost:6379"
	//if len(dsn) == 0 {
	//	dsn =
	//}
	client = redis.NewClient(&redis.Options{
		Addr: dsn,
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
}

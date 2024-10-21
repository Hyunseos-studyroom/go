package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

type Input struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8&parseTime=True&loc=Local"

	// 데이터베이스 연결 재시도 로직 추가
	var db *gorm.DB
	for {
		var err error
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Println("failed to connect to database, retrying in 5 seconds:", err)
		time.Sleep(5 * time.Second)
	}

	db.AutoMigrate(&Input{})

	r := gin.Default()

	r.POST("/post", func(c *gin.Context) {
		var i Input
		if err := c.ShouldBindJSON(&i); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		db.Create(&i)

		c.JSON(http.StatusCreated, gin.H{
			"input": i,
		})
	})

	r.GET("/get", func(c *gin.Context) {
		var inputs []Input
		db.Find(&inputs)

		c.JSON(http.StatusOK, gin.H{"data": inputs})
	})

	r.Run(os.Getenv("PORT"))
}

package main

import (
	"github.com/csv/controller"
	"github.com/csv/service"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	csvService := service.CSVService{}
	csvController := controller.Root{CSVService: csvService}

	r.POST("/upload", csvController.GetCSV)

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}

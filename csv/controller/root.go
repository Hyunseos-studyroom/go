package controller

import (
	"github.com/csv/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Root struct {
	CSVService service.CSVService
}

func (controller *Root) GetCSV(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "파일을 가져오는 데 실패했습니다."})
		return
	}

	data, err := controller.CSVService.ReadCSV(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CSV 파일을 읽는 데 실패했습니다."})
		return
	}

	c.JSON(http.StatusOK, data)
}

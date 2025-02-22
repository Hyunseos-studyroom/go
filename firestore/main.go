package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Firestore 클라이언트 초기화
func initFirestore() (*firestore.Client, error) {
	projectID := os.Getenv("FIRESTORE_PROJECT_ID")
	if projectID == "" {
		return nil, fmt.Errorf("FIRESTORE_PROJECT_ID environment variable not set")
	}

	// 로컬 개발용 인증 정보 (필요시 사용)
	// credFilePath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create firestore client: %v", err)
	}

	return client, nil
}

func main() {
	// Firestore 클라이언트 초기화
	fsClient, err := initFirestore()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize Firestore: %v", err))
	}

	r := gin.Default()
	
	// CORS 설정
	r.Use(cors.Default())

	r.POST("/upload", func(c *gin.Context) {
		// CSV 파일 가져오기
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{"error": "File upload error"})
			return
		}

		// 파일 열기
		src, err := file.Open()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to open file"})
			return
		}
		defer src.Close()

		// CSV 파싱
		reader := csv.NewReader(src)
		records, err := reader.ReadAll()
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid CSV format"})
			return
		}

		// 헤더와 데이터 분리
		if len(records) < 1 {
			c.JSON(400, gin.H{"error": "Empty CSV file"})
			return
		}
		headers := records[0]
		rows := records[1:]

		// Firestore 컬렉션 이름 가져오기
		collectionName := os.Getenv("FIRESTORE_COLLECTION")
		if collectionName == "" {
			c.JSON(500, gin.H{"error": "Collection name not configured"})
			return
		}

		// 데이터 저장
		ctx := context.Background()
		for _, row := range rows {
			data := make(map[string]interface{})
			for i, value := range row {
				if i < len(headers) {
					data[headers[i]] = value
				}
			}

			// Firestore에 문서 추가
			_, _, err := fsClient.Collection(collectionName).Add(ctx, data)
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to save data"})
				return
			}
		}

		c.JSON(200, gin.H{
			"message":  "Data saved successfully",
			"numRows":  len(rows),
		})
	})

	r.Run()
}

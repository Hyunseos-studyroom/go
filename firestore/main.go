package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	maxBatchSize         = 500              // Firestore batch write limit
	firestoreTimeout     = 5 * time.Minute  // Timeout for Firestore operations
	initFirestoreTimeout = 10 * time.Second // Timeout for Firestore client initialization
)

// Firestore 클라이언트 초기화
func initFirestore() (*firestore.Client, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("failed to load .env file: %v", err)
	}

	// 필수 환경 변수 검사.  main에서 해도 됨.
	projectID := os.Getenv("FIRESTORE_PROJECT_ID")
	if projectID == "" {
		return nil, fmt.Errorf("FIRESTORE_PROJECT_ID environment variable not set")
	}
	databaseID := os.Getenv("FIRESTORE_DATABASE_ID") // Optional, can be "(default)"
	credPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credPath == "" {
		return nil, fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS environment variable not set")
	}
	if os.Getenv("FIRESTORE_COLLECTION") == "" {
		return nil, fmt.Errorf("FIRESTORE_COLLECTION environment variable not set")
	}

	// 타임아웃 설정.  main에서 하는게 더 좋음.
	ctx, cancel := context.WithTimeout(context.Background(), initFirestoreTimeout)
	defer cancel()

	opts := []option.ClientOption{
		option.WithCredentialsFile(credPath),
		option.WithEndpoint("firestore.googleapis.com:443"),
	}
	// databaseID가 빈 문자열이면 기본값으로 설정됩니다.
	client, err := firestore.NewClientWithDatabase(ctx, projectID, databaseID, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create firestore client: %v", err)
	}

	return client, nil
}

// retryFirestoreOperation은 Firestore 작업에 대한 재시도 로직을 제공합니다.
func retryFirestoreOperation(ctx context.Context, operation func(ctx context.Context) error) error {
	backoff := 1 * time.Second
	maxBackoff := 30 * time.Second

	for attempt := 1; ; attempt++ {
		err := operation(ctx)
		if err == nil {
			return nil // Success
		}

		st, ok := status.FromError(err)
		if !ok {
			return err // Not a gRPC error, cannot retry
		}

		if st.Code() != codes.Unavailable && st.Code() != codes.DeadlineExceeded {
			return err // Not a retryable error
		}

		if attempt > 5 { // Max retries
			return fmt.Errorf("maximum retries exceeded: %w", err)
		}

		select {
		case <-ctx.Done():
			return ctx.Err() // Context cancelled or timed out
		case <-time.After(backoff): // Wait before retrying
			// Exponential backoff
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}

	}
}

func main() {
	// Firestore 클라이언트 초기화
	fsClient, err := initFirestore()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize Firestore: %v", err))
	}
	defer fsClient.Close()

	// 연결 테스트 (타임아웃 및 재시도 포함)
	ctx, cancel := context.WithTimeout(context.Background(), firestoreTimeout)
	defer cancel()

	err = retryFirestoreOperation(ctx, func(ctx context.Context) error {
		_, err := fsClient.Collections(ctx).GetAll() // collections 변수 사용 안함.
		return err
	})

	if err != nil {
		panic(fmt.Sprintf("Failed to list collections after retries: %v", err))
	}
	fmt.Printf("Successfully connected to Firestore.\n")

	r := gin.Default()

	// CORS 설정
	r.Use(cors.Default())

	// JSON 파일 업로드 및 Firestore로 가져오기
	r.POST("/import", func(c *gin.Context) {

		collectionName := os.Getenv("FIRESTORE_COLLECTION")

		// Multipart form data로 파일 받기
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
			return
		}

		// 파일 열기
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to open file: %v", err)})
			return
		}
		defer src.Close()

		// JSON 디코더 생성 (스트리밍 처리)
		decoder := json.NewDecoder(src)

		// 배열 토큰 '[' 읽기
		if _, err := decoder.Token(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format: expecting array"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), firestoreTimeout) // Gin 컨텍스트 사용 및 타임아웃 설정.
		defer cancel()

		var batch *firestore.WriteBatch
		batchSize := 0
		importedCount := 0

		// JSON 객체 스트리밍 및 Firestore 배치 처리
		for decoder.More() {
			var item map[string]interface{}
			if err := decoder.Decode(&item); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid JSON format: %v", err)})
				return
			}

			// 새 배치 시작 (필요한 경우)
			if batch == nil {
				batch = fsClient.Batch()
			}

			ref := fsClient.Collection(collectionName).NewDoc()
			batch.Set(ref, item)
			batchSize++
			importedCount++

			// 배치 크기 제한 도달 시 커밋
			if batchSize >= maxBatchSize {
				if _, err := batch.Commit(ctx); err != nil {

					c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to commit batch: %v", err)})
					return
				}
				batch = nil // 새 배치 준비
				batchSize = 0
			}

		}

		// 마지막 배치 커밋 (남은 데이터 처리)
		if batch != nil {
			if _, err := batch.Commit(ctx); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to commit final batch: %v", err)})
				return
			}
		}

		// 닫는 배열 토큰 ']' 읽기
		if _, err := decoder.Token(); err != nil && err != io.EOF { // io.EOF는 정상.
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format: expecting end of array"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Data imported successfully",
			"count":   importedCount,
		})
	})

	r.Run()
}

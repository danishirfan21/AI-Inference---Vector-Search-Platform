package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"github.com/your-repo/ai-platform/api/internal/handler"
	"github.com/your-repo/ai-platform/api/internal/repository"
	"github.com/your-repo/ai-platform/api/internal/service"
)

func main() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	natsSvc, err := service.NewNATSService(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "host=localhost user=user password=password dbname=ai_platform port=5432 sslmode=disable"
	}
	repo, err := repository.NewRepository(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}
	cacheSvc := service.NewCacheService(redisURL)

	milvusAddr := os.Getenv("MILVUS_ADDRESS")
	if milvusAddr == "" {
		milvusAddr = "localhost:19530"
	}
	milvusSvc, err := service.NewMilvusService(milvusAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Milvus: %v", err)
	}

	uploadHandler := handler.NewUploadHandler(natsSvc, repo)
	searchHandler := handler.NewSearchHandler(natsSvc, cacheSvc, milvusSvc)

	r := gin.Default()
	r.Use(otelgin.Middleware("api-service"))

	v1 := r.Group("/api/v1")
	{
		v1.POST("/upload", uploadHandler.Upload)
		v1.GET("/search", searchHandler.Search)
	}

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

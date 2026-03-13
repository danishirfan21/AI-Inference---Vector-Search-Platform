package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/your-repo/ai-platform/api/internal/model"
	"github.com/your-repo/ai-platform/api/internal/service"
)

import (
	"github.com/your-repo/ai-platform/api/internal/repository"
)

type UploadHandler struct {
	natsService *service.NATSService
	repo        *repository.Repository
}

func NewUploadHandler(nats *service.NATSService, repo *repository.Repository) *UploadHandler {
	return &UploadHandler{natsService: nats, repo: repo}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	contentType := c.PostForm("type")
	if contentType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Type is required"})
		return
	}

	id := uuid.New().String()
	s3Path := "uploads/" + id + "_" + file.Filename

	// In a real app: Save to S3 here

	content := &model.ContentMetadata{
		ID:               id,
		FileName:         file.Filename,
		FileType:         model.ContentType(contentType),
		S3Path:           s3Path,
		ProcessingStatus: model.StatusPending,
		CreatedAt:        time.Now(),
	}

	if err := h.repo.CreateContent(content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save metadata"})
		return
	}

	err = h.natsService.PublishContentUploaded(id, contentType, s3Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish event"})
		return
	}

	c.JSON(http.StatusAccepted, model.ContentMetadata{
		ID:               id,
		FileName:         file.Filename,
		FileType:         model.ContentType(contentType),
		S3Path:           s3Path,
		ProcessingStatus: model.StatusPending,
		CreatedAt:        time.Now(),
	})
}

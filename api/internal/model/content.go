package model

import (
	"time"
)

type ContentType string

const (
	ContentTypeDocument ContentType = "document"
	ContentTypeImage    ContentType = "image"
	ContentTypeAudio    ContentType = "audio"
)

type ProcessingStatus string

const (
	StatusPending    ProcessingStatus = "pending"
	StatusProcessing ProcessingStatus = "processing"
	StatusCompleted  ProcessingStatus = "completed"
	StatusFailed     ProcessingStatus = "failed"
)

type ContentMetadata struct {
	ID               string           `json:"id" gorm:"primaryKey"`
	FileName         string           `json:"file_name"`
	FileType         ContentType      `json:"file_type"`
	S3Path           string           `json:"s3_path"`
	ProcessingStatus ProcessingStatus `json:"processing_status"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}

type SearchResult struct {
	ID       string  `json:"id"`
	FileName string  `json:"file_name"`
	Score    float32 `json:"score"`
}

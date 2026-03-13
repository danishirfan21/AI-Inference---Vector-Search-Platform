package service

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

type NATSService struct {
	nc *nats.Conn
}

func NewNATSService(url string) (*NATSService, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}
	return &NATSService{nc: nc}, nil
}

func (s *NATSService) PublishContentUploaded(contentID string, fileType string, s3Path string) error {
	event := map[string]string{
		"content_id": contentID,
		"file_type":  fileType,
		"s3_path":    s3Path,
	}
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return s.nc.Publish("ContentUploaded", data)
}

func (s *NATSService) RequestEmbedding(query string) ([]float32, error) {
	msg, err := s.nc.Request("GetEmbedding", []byte(query), nats.DefaultTimeout)
	if err != nil {
		return nil, err
	}

	var embedding []float32
	if err := json.Unmarshal(msg.Data, &embedding); err != nil {
		return nil, err
	}
	return embedding, nil
}

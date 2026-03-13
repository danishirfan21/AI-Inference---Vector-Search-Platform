package service

import (
	"context"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type MilvusService struct {
	client client.Client
}

func NewMilvusService(addr string) (*MilvusService, error) {
	c, err := client.NewClient(context.Background(), client.Config{
		Address: addr,
	})
	if err != nil {
		return nil, err
	}
	return &MilvusService{client: c}, nil
}

func (s *MilvusService) Search(ctx context.Context, collectionName string, embedding []float32) ([]string, error) {
	sp, _ := entity.NewIndexIvfFlatSearchParam(10)
	searchResult, err := s.client.Search(
		ctx,
		collectionName,
		[]string{}, // partitions
		"",         // expr
		[]string{"id"},
		[]entity.Vector{entity.FloatVector(embedding)},
		"embedding",
		entity.L2,
		10,
		sp,
	)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, res := range searchResult {
		idColumn := res.Fields.GetColumn("id")
		if idColumn == nil {
			continue
		}
		for i := 0; i < idColumn.Len(); i++ {
			val, err := idColumn.GetAsString(i)
			if err == nil {
				ids = append(ids, val)
			}
		}
	}
	return ids, nil
}

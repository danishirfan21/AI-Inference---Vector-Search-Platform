package repository

import (
	"log"

	"github.com/your-repo/ai-platform/api/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.ContentMetadata{})
	if err != nil {
		log.Printf("AutoMigration failed: %v", err)
	}

	return &Repository{DB: db}, nil
}

func (r *Repository) CreateContent(content *model.ContentMetadata) error {
	return r.DB.Create(content).Error
}

func (r *Repository) GetContent(id string) (*model.ContentMetadata, error) {
	var content model.ContentMetadata
	err := r.DB.First(&content, "id = ?", id).Error
	return &content, err
}

package repository

import "gorm.io/gorm"

type AIModelRepository struct {
	db *gorm.DB
}

func NewAIModelRepository(db *gorm.DB) *AIModelRepository {
	return &AIModelRepository{
		db: db,
	}
}

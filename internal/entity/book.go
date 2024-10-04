package entity

import (
	"github.com/google/uuid"
	"os"
)

type Book struct {
	Id       string  `json:"id" gorm:"primaryKey;type:uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Title    string  `json:"title" example:"The Great Gatsby"`
	ISBN     string  `json:"isbn" example:"978-3-16-148410-0" gorm:"uniqueIndex"`
	AuthorId string  `json:"author_id" gorm:"uuid" example:"f47ac10b-58cc-4372-a567-0e02b2c3d479"`
	Author   *Author `json:"author,omitempty" gorm:"foreignKey:AuthorId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // Foreign key association with author
}

type UpsertBook struct {
	Title    string `json:"title" validate:"required,gte=2" example:"The Great Gatsby"`
	ISBN     string `json:"isbn" validate:"required,gte=2" example:"978-3-16-148410-0"`
	AuthorId string `json:"author_id" validate:"required,uuid" example:"f47ac10b-58cc-4372-a567-0e02b2c3d479"`
}

func (model *Book) GenerateModel(data *UpsertBook) {
	if model.Id == "" {
		model.Id = uuid.NewString()
	}
	model.Title = data.Title
	model.ISBN = data.ISBN
	model.AuthorId = data.AuthorId
}

func (model *Book) TableName() string {
	return os.Getenv("DB_PREFIX") + "book"
}

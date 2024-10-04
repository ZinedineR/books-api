package entity

import (
	"github.com/google/uuid"
	"os"
	"time"
)

type Author struct {
	Id        string     `json:"id" validate:"required,uuid" gorm:"primaryKey;type:uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name      string     `json:"name" validate:"required,title" example:"F. Scott Fitzgerald"`
	Birthdate *time.Time `json:"birthdate" example:"1896-09-24"`
}

type UpsertAuthor struct {
	Name      string `json:"name" validate:"required,gte=2" example:"Ernest Hemingway"`
	Birthdate string `json:"birthdate" validate:"required,dateLocal" example:"1899-07-21"`
}

func (model *Author) GenerateModel(data *UpsertAuthor) {
	if model.Id == "" {
		model.Id = uuid.NewString()
	}
	model.Name = data.Name
	birthDate, _ := time.Parse("2006-01-02", data.Birthdate)
	localBirthdate := birthDate.Local()
	model.Birthdate = &localBirthdate
}

func (model *Author) TableName() string {
	return os.Getenv("DB_PREFIX") + "author"
}

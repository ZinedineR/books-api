package entity

import (
	"os"
)

type User struct {
	Id       string `json:"id" gorm:"primaryKey;type:uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Username string `json:"username" example:"john_doe"`
	Password string `json:"password" example:"$2a$12$eixZaYVK1fsbw1ZfbX3OXe.PZyWJQ0Zf10hErsTQ6FVRHiA2vwLHu"` // Example of bcrypt-hashed password
}

type UserLogin struct {
	Username string `json:"username" validate:"required" example:"john_doe"`
	Password string `json:"password" validate:"required,password,gte=8" example:"SecurePass123!"` // "password" custom validation assumed
}

func (model *User) TableName() string {
	return os.Getenv("DB_PREFIX") + "user"
}

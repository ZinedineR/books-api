package entity

type Example struct {
	ID   int `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	Year int `json:"year" validate:"required"`
}

func (Example) TableName() string {
	return "example"
}

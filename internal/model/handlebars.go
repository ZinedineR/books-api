package model

type HandlebarsDetail struct {
	ID         string   `bson:"_id" json:"id" validate:"required,uuid"`
	Name       string   `bson:"name" json:"name" validate:"required"`
	Handlebars []string `json:"handlebars"`
}

package migration

import (
	"books-api/internal/entity"
	"books-api/pkg/database"
)

func AutoMigration(CpmDB *database.Database) {
	CpmDB.MigrateDB(

		&entity.User{},
		&entity.Author{},
		&entity.Book{})
	//&entity.SMSLog{}
}

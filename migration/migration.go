package migration

import (
	"boiler-plate-clean/internal/entity"
	"github.com/RumbiaID/pkg-library/app/pkg/database"
)

func AutoMigration(CpmDB *database.Database) {
	CpmDB.MigrateDB(

		&entity.Example{})
	//&entity.SMSLog{}
}

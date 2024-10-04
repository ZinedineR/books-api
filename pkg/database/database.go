package database

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"log/slog"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DbHost   string
	DbUser   string
	DbPass   string
	DbName   string
	DbPort   string
	DbPrefix string
}

type Database struct {
	db     *gorm.DB
	isCqrs bool
}

func (d *Database) GetDB() *gorm.DB {
	return d.db
}

func NewDatabase(driver string, cfg *Config) *Database {
	var db *gorm.DB
	var err error
	var dialect gorm.Dialector

	switch driver {
	case "postgres", "pgsql":
		dialect = postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", cfg.DbHost, cfg.DbUser, cfg.DbPass, cfg.DbName, cfg.DbPort))
	case "mysql":
		dialect = mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.DbUser, cfg.DbPass, cfg.DbHost, cfg.DbPort, cfg.DbName))
	case "sqlserver":
		dialect = sqlserver.Open(fmt.Sprintf("sqlserver://%s:%s@%s?database=%s", cfg.DbUser, cfg.DbPass, cfg.DbHost, cfg.DbName))
	case "oracle":
		slog.Warn("oracle driver is not supported yet")
		os.Exit(1)
	default:
		slog.Warn("unknown database driver")
		os.Exit(1)
	}

	for {
		configGorm := &gorm.Config{}
		if os.Getenv("APP_DEBUG") == "true" {
			configGorm.Logger = logger.Default.LogMode(logger.Info)
			// configGorm.DisableForeignKeyConstraintWhenMigrating = true
		}
		configGorm.NamingStrategy = schema.NamingStrategy{
			TablePrefix: cfg.DbPrefix, // table name prefix, table for `User` would be `t_users`
		}
		db, err = gorm.Open(dialect, configGorm)
		if err != nil {
			slog.Error(fmt.Sprintf("failed to connect to %s database", driver), "error", err.Error())
			slog.Info(fmt.Sprintf("retrying to connect to %s database in 5 seconds...", driver))
			time.Sleep(5 * time.Second)
			continue
		}
		slog.Info(fmt.Sprintf("successfully connected to %s database", driver))
		break
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("failed to configure connection pool", "error", err.Error())
		os.Exit(1)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &Database{db: db, isCqrs: false}
}

func (m *Database) CqrsDB(driver string, cfg *Config) {
	var dialectRead gorm.Dialector

	switch driver {
	case "postgres", "pgsql":
		dialectRead = postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", cfg.DbHost, cfg.DbUser, cfg.DbPass, cfg.DbName, cfg.DbPort))
	case "mysql":
		dialectRead = mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.DbUser, cfg.DbPass, cfg.DbHost, cfg.DbPort, cfg.DbName))
	case "sqlserver":
		dialectRead = sqlserver.Open(fmt.Sprintf("sqlserver://%s:%s@%s?database=%s", cfg.DbUser, cfg.DbPass, cfg.DbHost, cfg.DbName))
	case "oracle":
		slog.Warn("oracle driver is not supported yet")
		os.Exit(1)
	default:
		slog.Warn("unknown database driver")
		os.Exit(1)
	}

	// Read DB
	resolver := dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{
			dialectRead,
		},
		TraceResolverMode: true,
	})
	errResolver := m.db.Use(resolver)
	if errResolver != nil {
		slog.Error("failed to configure connection pool READ", "error", errResolver.Error())
		os.Exit(1)
	}
	//change isCqrs status
	m.isCqrs = true
}

func (m *Database) MigrateDB(dst ...interface{}) {

	if m.isCqrs {

		// CQRS
		err := m.db.Clauses(dbresolver.Write).AutoMigrate(dst...)
		if err != nil {
			slog.Error("failed to migrate db", "error", err.Error())
			os.Exit(1)
		}

		err = m.db.Clauses(dbresolver.Read).AutoMigrate(dst...)
		if err != nil {
			slog.Error("failed to migrate db Read", "error", err.Error())
			os.Exit(1)
		}
	} else {

		// Not CQRS
		err := m.db.AutoMigrate(dst...)
		if err != nil {
			slog.Error("failed to migrate db", "error", err.Error())
			os.Exit(1)
		}
	}

	slog.Info(fmt.Sprintf("successfully migrated entity: %T", dst...))
}

func (m *Database) DownMigrate(all bool, dst ...interface{}) {
	if all {
		if m.isCqrs {

			// CQRS
			err := m.db.Clauses(dbresolver.Read).Migrator().DropTable(dst)
			if err != nil {
				slog.Error("failed to drop index db", "error", err.Error())
			}

			err = m.db.Clauses(dbresolver.Write).Migrator().DropTable(dst)
			if err != nil {
				slog.Error("failed to drop index db", "error", err.Error())
			}

		} else {
			err := m.db.Migrator().DropTable(dst)
			if err != nil {
				slog.Error("failed to drop index db", "error", err.Error())
			}
		}

	} else if !all {
		for _, target := range dst {
			m.WipeTable(target)
		}
	}
	slog.Info(fmt.Sprintf("successfully migrated entity: %T", dst...))
}

func (m *Database) DropColumnDB(dst interface{}, columnTarget string) {
	if m.isCqrs {

		// CQRS
		err := m.db.Clauses(dbresolver.Read).Migrator().DropColumn(dst, columnTarget)
		if err != nil {
			slog.Error("failed to delete column", "error", err.Error())
		}

		err = m.db.Clauses(dbresolver.Write).Migrator().DropColumn(dst, columnTarget)
		if err != nil {
			slog.Error("failed to delete column", "error", err.Error())
		}
	} else {

		// Not CQRS
		err := m.db.Migrator().DropColumn(dst, columnTarget)
		if err != nil {
			slog.Error("failed to delete column", "error", err.Error())
		}
	}

	slog.Info(fmt.Sprintf("successfully migrated entity: %T", dst))
}

func (m *Database) RenameColumnDB(dst interface{}, oldname, columnTarget string) {

	if m.isCqrs {

		// CQRS
		err := m.db.Clauses(dbresolver.Write).Migrator().RenameColumn(dst, oldname, columnTarget)
		if err != nil {
			slog.Error("failed to rename column", "error", err.Error())
		}

		err = m.db.Clauses(dbresolver.Read).Migrator().RenameColumn(dst, oldname, columnTarget)
		if err != nil {
			slog.Error("failed to rename column", "error", err.Error())
		}
	} else {
		err := m.db.Migrator().RenameColumn(dst, oldname, columnTarget)
		if err != nil {
			slog.Error("failed to rename column", "error", err.Error())
		}
	}
	slog.Info(fmt.Sprintf("successfully migrated entity: %T", dst))
}

func (m *Database) DownIndexDB(dst interface{}, columnTarget string) {

	if m.isCqrs {

		// CQRS
		index, err := m.db.Clauses(dbresolver.Read).Migrator().GetIndexes(dst)
		for _, indexData := range index {
			column := indexData.Columns()
			for _, columnName := range column {
				if columnName == columnTarget {
					err = m.db.Clauses(dbresolver.Read).Migrator().DropIndex(dst, indexData.Name())
					if err = m.db.Clauses(dbresolver.Read).Migrator().DropConstraint(dst, indexData.Name()); err != nil {
						logrus.Error("failed to drop index db", "error", err.Error())
						//os.Exit(1)
					}
				}
			}
		}
		if err != nil {
			slog.Error("failed to drop index db", "error", err.Error())
			//os.Exit(1)
		}

		index, err = m.db.Clauses(dbresolver.Write).Migrator().GetIndexes(dst)
		for _, indexData := range index {
			column := indexData.Columns()
			for _, columnName := range column {
				if columnName == columnTarget {
					err = m.db.Clauses(dbresolver.Write).Migrator().DropIndex(dst, indexData.Name())
					if err = m.db.Clauses(dbresolver.Write).Migrator().DropConstraint(dst, indexData.Name()); err != nil {
						logrus.Error("failed to drop index db", "error", err.Error())
						//os.Exit(1)
					}
				}
			}
		}

		if err != nil {
			slog.Error("failed to drop index db", "error", err.Error())
			//os.Exit(1)
		}

	} else {

		// NOT CQRS
		index, err := m.db.Migrator().GetIndexes(dst)
		for _, indexData := range index {
			column := indexData.Columns()
			for _, columnName := range column {
				if columnName == columnTarget {
					err = m.db.Migrator().DropIndex(dst, indexData.Name())
					if err = m.db.Migrator().DropConstraint(dst, indexData.Name()); err != nil {
						logrus.Error("failed to drop index db", "error", err.Error())
						//os.Exit(1)
					}
				}
			}
		}

		if err != nil {
			slog.Error("failed to drop index db", "error", err.Error())
			//os.Exit(1)
		}
	}

	slog.Info(fmt.Sprintf("successfully migrated entity: %T", dst))
}

func (m *Database) WipeTable(dst interface{}) {
	if m.isCqrs {

		// CQRS
		err := m.db.Clauses(dbresolver.Read).Where("id is not null").Delete(dst).Error
		if err != nil {
			slog.Error("failed to drop index db", "error", err.Error())
		}

		err = m.db.Clauses(dbresolver.Write).Where("id is not null").Delete(dst).Error
		if err != nil {
			slog.Error("failed to drop index db", "error", err.Error())
		}
	} else {

		// NOT CQRS
		err := m.db.Where("id is not null").Delete(dst).Error
		if err != nil {
			slog.Error("failed to drop index db", "error", err.Error())
		}
	}

}

func (m *Database) DeleteTable(dst ...interface{}) {
	if m.isCqrs {

		// CQRS
		err := m.db.Clauses(dbresolver.Read).Migrator().DropTable(dst)
		if err != nil {
			slog.Error("failed to drop index db", "error", err.Error())
		}

		err = m.db.Clauses(dbresolver.Write).Migrator().DropTable(dst)
		if err != nil {
			slog.Error("failed to drop index db", "error", err.Error())
		}

	} else {

		// Not CQRS
		err := m.db.Migrator().DropTable(dst)
		if err != nil {
			slog.Error("failed to drop index db", "error", err.Error())
		}
	}
}

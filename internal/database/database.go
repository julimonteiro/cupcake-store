package database

import (
	"fmt"
	"log"

	"github.com/julimonteiro/cupcake-store/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(cfg *config.Config) (db *gorm.DB, err error) {
	gormLogger := logger.Default.LogMode(logger.Info)
	if cfg.LogLevel == "error" {
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	switch cfg.DBDialect {
	case "postgres":
		db, err = gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{
			Logger: gormLogger,
		})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(cfg.DBDSN), &gorm.Config{
			Logger: gormLogger,
		})
	default:
		return nil, fmt.Errorf("unsupported database dialect: %s", cfg.DBDialect)
	}

	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("error running migrations: %w", err)
	}

	log.Printf("Connected to database %s", cfg.DBDialect)
	return db, nil
}

func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
	// TODO
	)
}
